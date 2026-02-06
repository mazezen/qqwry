import requests
import os

url = "https://github.com/metowolf/qqwry.dat/releases/latest/download/qqwry.dat"

response = requests.get(url)

if response.status_code == 200 and len(response.content) > 0:
    file_path = os.path.join(os.path.dirname(__file__), "qqwry.dat")
    with open(file_path, "wb") as file:
        file.write(response.content)
    file_size = os.path.getsize("qqwry.dat") / (1024 * 1024)  # 转换为 MB
    print(f"下载成功！文件大小：{file_size:.2f} MB")
else:
    print("下载失败，请检查网络连接或稍后重试。")