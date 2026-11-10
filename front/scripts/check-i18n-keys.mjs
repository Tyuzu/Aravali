import fs from 'node:fs';

const files = [
  'public/i18n/en.json',
  'public/i18n/es.json',
  'public/i18n/fr.json',
  'public/i18n/hi.json',
  'public/i18n/ar.json',
  'public/i18n/jp.json'
];

const checks = [
  'nav.baitos',
  'nav.more',
  'footer.copyright',
  'footer.lang.en',
  'nav.dash',
  'nav.tools'
];

let failed = false;

for (const file of files) {
  const data = JSON.parse(fs.readFileSync(file, 'utf8'));

  for (const key of checks) {
    const segments = key.split('.');
    let value = data;

    for (let index = 0; index < segments.length; index += 1) {
      const segment = segments[index];

      if (value == null) {
        value = undefined;
        break;
      }

      if (typeof value !== 'object') {
        value = undefined;
        break;
      }

      const record = value;
      const exactMatch = record[segment];

      if (Object.prototype.hasOwnProperty.call(record, segment)) {
        value = exactMatch;
        continue;
      }

      const remainder = segments.slice(index).join('.');
      if (Object.prototype.hasOwnProperty.call(record, remainder)) {
        value = record[remainder];
        break;
      }

      value = undefined;
      break;
    }

    if (typeof value !== 'string') {
      console.log('MISSING', file, key);
      failed = true;
    }
  }
}

if (failed) {
  console.log('MISSING_KEYS');
  process.exit(1);
}

console.log('ALL_KEYS_PRESENT');
