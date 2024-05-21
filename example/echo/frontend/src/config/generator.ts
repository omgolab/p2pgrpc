import {getImage} from 'astro:assets';
import {log} from 'console';
import fs from 'fs';
import path from 'path';
import {vars} from './variables';

export async function transformRemoteImage(
  imgUrl: string,
  destPath: string,
  format: 'webp' | 'png' | 'avif' = 'png',
  height = 512,
  width = 512
): Promise<void> {
  const res = await getImage({
    src: imgUrl,
    format: format,
    height: height,
    width: width,
  });
  const sourcePath = path.resolve(
    vars.OUT_DIR,
    res.src.replace(vars.env.APP_ROOT, '')
  );
  const destinationPath = path.resolve(destPath.replace(vars.env.APP_ROOT, ''));
  const mk = fs.promises.mkdir(path.dirname(destinationPath), {
    recursive: true,
  });

  let fileExists = fs.existsSync(sourcePath);
  while (!fileExists) {
    await new Promise(resolve => setTimeout(resolve, 100));
    fileExists = fs.existsSync(sourcePath);
  }

  await mk;
  log(`Moving ${sourcePath} to ${destinationPath}`);
  fs.promises.rename(sourcePath, destinationPath).catch(err => {
    log(err);
  });
}
