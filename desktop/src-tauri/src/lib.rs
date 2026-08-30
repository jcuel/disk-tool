use std::sync::Mutex;
use std::time::{Duration, Instant};

use tauri::webview::WebviewWindowBuilder;
use tauri::{Manager, RunEvent, Url, WebviewUrl};
use tauri_plugin_shell::process::{CommandChild, CommandEvent};
use tauri_plugin_shell::ShellExt;
use tokio::time::timeout;

const READY_PREFIX: &str = "disk-tool-ready port=";
const STARTUP_TIMEOUT: Duration = Duration::from_secs(30);

struct SidecarState {
    child: Mutex<Option<CommandChild>>,
}

fn parse_ready_port(line: &str) -> Option<u16> {
    let line = line.trim();
    line
        .strip_prefix(READY_PREFIX)?
        .parse()
        .ok()
        .filter(|p| *p > 0)
}

fn allows_localhost(url: &Url) -> bool {
    url.scheme() == "about"
        || (url.scheme() == "http"
            && matches!(url.host_str(), Some("127.0.0.1" | "localhost")))
}

async fn start_sidecar(app: &tauri::AppHandle) -> Result<u16, String> {
    let sidecar = app
        .shell()
        .sidecar("disk-tool")
        .map_err(|e| format!("sidecar command: {e}"))?
        .args(["serve", "--no-open", "--port", "0", "--ready-stdout"]);

    let (mut rx, child) = sidecar
        .spawn()
        .map_err(|e| format!("spawn sidecar: {e}"))?;

    let deadline = Instant::now() + STARTUP_TIMEOUT;
    let mut buffer = String::new();
    let mut port: Option<u16> = None;
    'wait: while port.is_none() {
        if Instant::now() >= deadline {
            let _ = child.kill();
            return Err("sidecar startup timed out".into());
        }

        let remaining = deadline.saturating_duration_since(Instant::now());
        let event = match timeout(remaining, rx.recv()).await {
            Ok(Some(ev)) => ev,
            Ok(None) => {
                let _ = child.kill();
                return Err("sidecar exited before ready".into());
            }
            Err(_) => {
                let _ = child.kill();
                return Err("sidecar startup timed out".into());
            }
        };

        match event {
            CommandEvent::Stdout(bytes) => {
                buffer.push_str(&String::from_utf8_lossy(&bytes));
                while let Some(pos) = buffer.find('\n') {
                    let line = buffer[..pos].to_string();
                    buffer.drain(..=pos);
                    if let Some(p) = parse_ready_port(&line) {
                        port = Some(p);
                        break 'wait;
                    }
                }
            }
            CommandEvent::Terminated(payload) => {
                let _ = child.kill();
                return Err(format!(
                    "sidecar exited before ready (code={:?}, signal={:?})",
                    payload.code, payload.signal
                ));
            }
            CommandEvent::Error(err) => {
                let _ = child.kill();
                return Err(format!("sidecar error: {err}"));
            }
            _ => {}
        }
    }
    let port = port.expect("ready port set");

    if let Some(state) = app.try_state::<SidecarState>() {
        *state.child.lock().unwrap() = Some(child);
    }

    Ok(port)
}

fn kill_sidecar(app: &tauri::AppHandle) {
    if let Some(state) = app.try_state::<SidecarState>() {
        if let Some(child) = state.child.lock().unwrap().take() {
            let _ = child.kill();
        }
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .manage(SidecarState {
            child: Mutex::new(None),
        })
        .setup(|app| {
            let handle = app.handle().clone();
            tauri::async_runtime::block_on(async move {
                let port = start_sidecar(&handle).await?;
                let url = format!("http://127.0.0.1:{port}");
                let parsed =
                    Url::parse(&url).map_err(|e| format!("invalid sidecar url: {e}"))?;

                WebviewWindowBuilder::new(&handle, "main", WebviewUrl::External(parsed))
                    .title("disk-tool")
                    .inner_size(1280.0, 800.0)
                    .resizable(true)
                    .on_navigation(|url| allows_localhost(url))
                    .build()
                    .map_err(|e| format!("create main window: {e}"))?;

                Ok::<(), String>(())
            })
            .map_err(|e| -> Box<dyn std::error::Error> { e.into() })?;
            Ok(())
        })
        .build(tauri::generate_context!())
        .expect("error while building tauri application")
        .run(|app, event| {
            if matches!(event, RunEvent::Exit) {
                kill_sidecar(app);
            }
        });
}
