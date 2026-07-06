#!/usr/bin/env python3
"""Generate landing page PNG screenshots and og-image. Run: python scripts/generate-landing-images.py"""
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

ROOT = Path(__file__).resolve().parent.parent
OUT = ROOT / "site" / "public" / "images"
BG = (15, 20, 25)
PANEL = (22, 27, 34)
BORDER = (48, 54, 61)
TEXT = (230, 237, 243)
MUTED = (139, 148, 158)
ACCENT = (35, 134, 54)
BLUE = (56, 139, 253)


def font(size: int, bold: bool = False):
    try:
        name = "arialbd.ttf" if bold else "arial.ttf"
        return ImageFont.truetype(name, size)
    except OSError:
        return ImageFont.load_default()


def draw_overview(path: Path) -> None:
    w, h = 800, 450
    img = Image.new("RGB", (w, h), BG)
    d = ImageDraw.Draw(img)
    d.rectangle([0, 0, w, 56], fill=PANEL, outline=BORDER)
    d.text((24, 16), "disk-tool", fill=TEXT, font=font(18, True))
    d.text((24, 68), "/demo/projects", fill=MUTED, font=font(12))
    d.rectangle([24, 92, 384, 422], fill=PANEL, outline=BORDER)
    d.text((40, 108), "Folder tree", fill=TEXT, font=font(14, True))
    d.text((48, 140), "big-dir", fill=BLUE, font=font(13))
    d.text((280, 140), "199 B", fill=MUTED, font=font(12))
    d.text((48, 168), "small-dir", fill=TEXT, font=font(13))
    d.text((280, 168), "6 B", fill=MUTED, font=font(12))
    d.rectangle([404, 92, 776, 252], fill=PANEL, outline=BORDER)
    d.text((420, 108), "Distribution", fill=TEXT, font=font(14, True))
    d.rectangle([440, 140, 560, 230], fill=ACCENT)
    d.rectangle([570, 160, 630, 230], fill=BLUE)
    d.rectangle([404, 268, 776, 422], fill=PANEL, outline=BORDER)
    d.text((420, 284), "Insights", fill=TEXT, font=font(14, True))
    d.text((420, 312), "big-dir uses 90% of scanned space", fill=MUTED, font=font(12))
    img.save(path, optimize=True)


def draw_insights(path: Path) -> None:
    w, h = 800, 450
    img = Image.new("RGB", (w, h), BG)
    d = ImageDraw.Draw(img)
    d.rectangle([24, 24, 776, 426], fill=PANEL, outline=BORDER)
    d.text((40, 44), "Insights and cleanup", fill=TEXT, font=font(16, True))
    d.rectangle([40, 76, 760, 124], fill=BG, outline=BORDER)
    d.text((56, 94), "Maintenance presets - Dev reclaim - Temp cleanup", fill=MUTED, font=font(13))
    d.rectangle([40, 140, 760, 204], fill=(31, 61, 42), outline=ACCENT)
    d.text((56, 158), "review", fill=TEXT, font=font(13))
    d.text((56, 178), "node_modules - regenerable dev artifact", fill=MUTED, font=font(12))
    d.rectangle([40, 220, 760, 284], fill=BG, outline=BORDER)
    d.text((56, 238), "caution", fill=TEXT, font=font(13))
    d.text((56, 258), ".npm cache - review before delete", fill=MUTED, font=font(12))
    d.rectangle([40, 310, 200, 346], fill=ACCENT)
    d.text((72, 322), "Copy report", fill=(255, 255, 255), font=font(13, True))
    d.rectangle([212, 310, 332, 346], fill=PANEL, outline=BORDER)
    d.text((244, 322), "Support ticket", fill=TEXT, font=font(13))
    img.save(path, optimize=True)


def draw_og(path: Path) -> None:
    w, h = 1200, 630
    img = Image.new("RGB", (w, h), BG)
    d = ImageDraw.Draw(img)
    d.text((60, 180), "disk-tool", fill=TEXT, font=font(72, True))
    d.text(
        (60, 280),
        "Find where disk space goes",
        fill=MUTED,
        font=font(36),
    )
    d.text(
        (60, 340),
        "Local overview-first scan  |  Cleanup insights  |  Try the demo",
        fill=MUTED,
        font=font(24),
    )
    d.rectangle([60, 420, 320, 480], fill=ACCENT)
    d.text((100, 438), "Download", fill=(255, 255, 255), font=font(28, True))
    img.save(path, optimize=True)


def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    draw_overview(OUT / "overview.png")
    draw_insights(OUT / "insights.png")
    draw_og(ROOT / "site" / "public" / "og-image.png")
    print(f"Wrote images to {OUT} and og-image.png")


if __name__ == "__main__":
    main()
