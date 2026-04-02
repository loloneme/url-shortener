import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:8080'

export const options = {
    stages: [
        {duration: '10s', target: 10},
        {duration: '1m', target: 50},
        {duration: '2m', target: 100},
        {duration: '3m', target: 500},
        {duration: '5m', target: 1000},
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'], 
        http_req_duration: ['p(95)<200'],
      },
};

function shortenUrl(original){
    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'ShortenUrl' },
        timeout: '10s',
    };

    const payload = JSON.stringify({ url: original });
    const res = http.post(`${BASE_URL}/api/shorten`, payload, params);

    const success = check(res, {
        'shorten: status 200 or 201': (r) => r.status === 200 || r.status === 201,
        'shorten: has short_url':  (r) => !!r.json('short_url'),
    });

    if (success) {
        return res.json('short_url');
    }

    return null;
}

function getOriginalUrl(short){
    const params = {
        tags: { name: 'GetOriginalUrl' },
        timeout: '10s',
    };

    const res = http.get(`${BASE_URL}/api/${short}`, params);

    const success = check(res, {
        'get original url: status 200': (r) => r.status === 200,
        'get original url: has original_url': (r) => !!r.json('original_url'),
    });

}


export default function () {
    const originalUrl = `https://example.com/${__VU}-${__ITER}`;

    const shortUrl = shortenUrl(originalUrl);

    const short = shortUrl.split('/')[3];

    sleep(0.5)

    if (shortUrl) {
        getOriginalUrl(short);
    }

    sleep(0.5)
}