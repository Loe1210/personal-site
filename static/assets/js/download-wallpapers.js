const https = require('https');
const fs = require('fs');
const path = require('path');

const WALLPAPER_DIR = path.join(__dirname, '..', 'wallpapers');
const OUTPUT_JS = path.join(__dirname, '..', 'json', 'images.js');

// 获取前20张 Bing 壁纸（idx从0到16，每次取8张）
async function fetchBingImages() {
    const allImages = [];
    // 从多个市场获取壁纸（不同市场的壁纸不完全一样）
    const markets = [
        { host: 'cn.bing.com', mkt: 'zh-CN' },
        { host: 'www.bing.com', mkt: 'en-US' },
        { host: 'www.bing.com', mkt: 'en-GB' },
        { host: 'www.bing.com', mkt: 'ja-JP' },
        { host: 'www.bing.com', mkt: 'de-DE' },
        { host: 'www.bing.com', mkt: 'fr-FR' },
    ];
    for (const { host, mkt } of markets) {
        for (let idx = 0; idx < 3; idx++) {
            const url = `https://${host}/HPImageArchive.aspx?format=js&idx=${idx}&n=8&mkt=${mkt}`;
            try {
                const data = await new Promise((resolve, reject) => {
                    https.get(url, (res) => {
                        let body = '';
                        res.on('data', chunk => body += chunk);
                        res.on('end', () => {
                            try { resolve(JSON.parse(body)); } catch (e) { reject(e); }
                        });
                    }).on('error', reject);
                });
                if (data && data.images) {
                    data.images.forEach(img => allImages.push(img));
                }
            } catch (e) {
                console.error(`Failed to fetch host=${host} mkt=${mkt} idx=${idx}:`, e.message);
            }
        }
    }
    // 去重并取前15张
    const unique = [];
    const seen = new Set();
    for (const img of allImages) {
        if (!seen.has(img.urlbase) && unique.length < 15) {
            seen.add(img.urlbase);
            unique.push(img);
        }
    }
    return unique;
}

// 下载单张图片
function downloadImage(urlbase, filename) {
    return new Promise((resolve, reject) => {
        const url = `https://www.bing.com${urlbase}_1920x1080.jpg`;
        const filePath = path.join(WALLPAPER_DIR, filename);
        const file = fs.createWriteStream(filePath);
        https.get(url, (res) => {
            res.pipe(file);
            file.on('finish', () => {
                file.close();
                console.log(`Downloaded: ${filename}`);
                resolve();
            });
        }).on('error', (err) => {
            fs.unlink(filePath, () => {});
            reject(err);
        });
    });
}

async function main() {
    console.log('Fetching Bing images...');
    const images = await fetchBingImages();
    console.log(`Found ${images.length} images`);

    const localFiles = [];
    for (let i = 0; i < images.length; i++) {
        const filename = `bing-${String(i + 1).padStart(2, '0')}.jpg`;
        await downloadImage(images[i].urlbase, filename);
        localFiles.push(`wallpapers/${filename}`);
    }

    // 生成 images.js 指向本地文件
    const jsContent = `window.BING_IMAGES = ${JSON.stringify(localFiles, null, 2)};\n`;
    fs.writeFileSync(OUTPUT_JS, jsContent, 'utf-8');
    console.log(`\nDone! Saved ${localFiles.length} wallpapers to ${WALLPAPER_DIR}`);
    console.log(`Updated ${OUTPUT_JS}`);
}

main().catch(err => {
    console.error('Error:', err);
    process.exit(1);
});