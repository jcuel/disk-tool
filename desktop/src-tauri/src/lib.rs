use std::sync::Mutex;
use std::time::{Duration, Instant};

use tauri::{Manager, RunEvent, Url};
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
        let event = timeout(remaining, rx.recv())
            .await
            .map_err(|_| "sidecar startup timed out".to_string())?
            .ok_or_else(|| "sidecar exited before ready".to_string())?;

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
                return Err(format!(
                    "sidecar exited before ready (code={:?}, signal={:?})",
                    payload.code, payload.signal
                ));
            }
            CommandEvent::Error(err) => {
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
                let window = handle
                    .get_webview_window("main")
                    .ok_or("main window not found")?;
                window
                    .navigate(
                        Url::parse(&url).map_err(|e| format!("invalid sidecar url: {e}"))?,
                    )
                    .map_err(|e| format!("navigate webview: {e}"))?;
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
