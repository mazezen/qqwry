<?php
// PHP 纯真 IP 地址数据库自动更新功能
$url = "https://github.com/metowolf/qqwry.dat/releases/latest/download/qqwry.dat";
$dat = file_get_contents($url);

if ($dat === false || strlen($dat) < 100000) {
    die("下载失败或文件过小，更新终止。");
}

$fp = fopen("qqwry.dat", "wb");
if ($fp) {
    fwrite($fp, $dat);
    fclose($fp);
    echo "更新成功！文件大小：" . round(filesize("qqwry.dat") / 1024 / 1024, 2) . " MB";
} else {
    echo "写入失败，请检查目录权限";
}


?>