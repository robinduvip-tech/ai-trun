import os
import math
from PIL import Image, ImageDraw

def draw_logo(size):
    # 创建透明底图
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # 缩放比率 (以 512 像素为基准)
    scale = size / 512.0
    cx, cy = size / 2.0, size / 2.0

    # 渐变配色方案 (Blue -> Indigo -> Emerald Green)
    color_blue = (59, 130, 246, 255)
    color_indigo = (99, 102, 241, 255)
    color_emerald = (16, 185, 129, 255)

    def get_gradient_color(t):
        # 线性插值计算流光渐变
        if t <= 0.5:
            factor = t / 0.5
            r = int(color_blue[0] + (color_indigo[0] - color_blue[0]) * factor)
            g = int(color_blue[1] + (color_indigo[1] - color_blue[1]) * factor)
            b = int(color_blue[2] + (color_indigo[2] - color_blue[2]) * factor)
        else:
            factor = (t - 0.5) / 0.5
            r = int(color_indigo[0] + (color_emerald[0] - color_indigo[0]) * factor)
            g = int(color_indigo[1] + (color_emerald[1] - color_indigo[1]) * factor)
            b = int(color_indigo[2] + (color_emerald[2] - color_indigo[2]) * factor)
        return (r, g, b, 255)

    # 1. 绘制外部虚线轨道 (r=210)
    orbit_r = 210 * scale
    num_segments = 16
    for i in range(num_segments):
        start_angle = i * (360 / num_segments)
        end_angle = start_angle + 12
        # 画弧段
        for angle_deg in range(int(start_angle), int(end_angle)):
            angle_rad = math.radians(angle_deg)
            x = cx + orbit_r * math.cos(angle_rad)
            y = cy + orbit_r * math.sin(angle_rad)

            # 渐变吸收
            t = (angle_deg % 360) / 360.0
            col = get_gradient_color(t)
            # 外部虚线轨道具有 0.6 的透明度
            col_with_alpha = (col[0], col[1], col[2], int(255 * 0.65))

            dot_r = 4.5 * scale if size > 128 else 1.0 * scale
            draw.ellipse([x - dot_r, y - dot_r, x + dot_r, y + dot_r], fill=col_with_alpha)

    # 2. 绘制 "C" 字型弧线 (左翼)
    c_r = 140 * scale
    stroke_w = int(24 * scale) if size > 64 else max(3, int(16 * scale))

    # 采用密集画圆点的方式实现具有完美顺滑渐变的 C 弧线
    for deg in range(125, 235):
        rad = math.radians(deg)
        x = cx + c_r * math.cos(rad)
        y = cy + c_r * math.sin(rad)
        t = deg / 360.0
        col = get_gradient_color(t)
        draw.ellipse([x - stroke_w/2, y - stroke_w/2, x + stroke_w/2, y + stroke_w/2], fill=col)

    # 3. 绘制 "X" 字型右翼路由交叉
    # 右翼从 (410, 140) 到中心 (256, 256)，再到 (410, 372)
    # 采用等分段插值，实现具有完美色彩过渡的渐变实线
    steps = 150
    for s in range(steps):
        factor = s / float(steps)
        # 上射束
        x1 = cx + (154 * scale) * factor
        y1 = cy - (124 * scale) * factor
        t1 = 0.5 + (0.3 * factor)
        draw.ellipse([x1 - stroke_w/2, y1 - stroke_w/2, x1 + stroke_w/2, y1 + stroke_w/2], fill=get_gradient_color(t1))

        # 下射束
        x2 = cx + (154 * scale) * factor
        y2 = cy + (124 * scale) * factor
        t2 = 0.5 + (0.3 * factor)
        draw.ellipse([x2 - stroke_w/2, y2 - stroke_w/2, x2 + stroke_w/2, y2 + stroke_w/2], fill=get_gradient_color(t2))

    # 4. 绘制 X 的左侧反向贯穿路径 (从中心 256,256 到 170,170 和 170,342)
    for s in range(steps):
        factor = s / float(steps)
        # 左上射束
        x1 = cx - (100 * scale) * factor
        y1 = cy - (100 * scale) * factor
        t1 = 0.5 - (0.3 * factor)
        draw.ellipse([x1 - stroke_w/2, y1 - stroke_w/2, x1 + stroke_w/2, y1 + stroke_w/2], fill=get_gradient_color(t1))

        # 左下射束
        x2 = cx - (100 * scale) * factor
        y2 = cy + (100 * scale) * factor
        t2 = 0.5 - (0.3 * factor)
        draw.ellipse([x2 - stroke_w/2, y2 - stroke_w/2, x2 + stroke_w/2, y2 + stroke_w/2], fill=get_gradient_color(t2))

    # 5. 绘制核心 AI 路由核 (Center Glowing Node)
    core_glow_r = 45 * scale
    # 渐变发光外层晕 (0.5 opacity)
    for r_offset in range(int(core_glow_r), 0, -1):
        alpha_factor = 0.5 * (1.0 - (r_offset / core_glow_r))
        col_alpha = get_gradient_color(0.5)
        fill_col = (col_alpha[0], col_alpha[1], col_alpha[2], int(255 * alpha_factor))
        draw.ellipse([cx - r_offset, cy - r_offset, cx + r_offset, cy + r_offset], fill=fill_col)

    # 内层纯白实心核心
    white_core_r = 24 * scale
    draw.ellipse([cx - white_core_r, cy - white_core_r, cx + white_core_r, cy + white_core_r], fill=(255, 255, 255, 255))

    return img

def main():
    print("🎨 开始生成 CCX 高清品牌图标资产...")

    # A. 替换桌面打包主图标 desktop/build/appicon.png (512x512)
    appicon_path = "desktop/build/appicon.png"
    appicon_dir = os.path.dirname(appicon_path)
    if not os.path.exists(appicon_dir):
        os.makedirs(appicon_dir, exist_ok=True)
    img_512 = draw_logo(512)
    img_512.save(appicon_path, "PNG")
    print(f"✅ 生成成功：{appicon_path} (512x512 高清)")

    # B. 替换托盘图标源 desktop/frontend/public/wails.png (32x32)
    wails_path = "desktop/frontend/public/wails.png"
    wails_dir = os.path.dirname(wails_path)
    if not os.path.exists(wails_dir):
        os.makedirs(wails_dir, exist_ok=True)
    img_32 = draw_logo(32)
    img_32.save(wails_path, "PNG")
    print(f"✅ 生成成功：{wails_path} (32x32 托盘)")

    print("🎉 官方品牌图标资产全部替换、覆盖完毕！")

if __name__ == "__main__":
    main()
