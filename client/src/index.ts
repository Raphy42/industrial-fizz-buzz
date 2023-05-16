import {delay} from 'https://deno.land/std@0.187.0/async/delay.ts';

const apiUrl = Deno.env.get('API_URL');
if (apiUrl === undefined) {
    throw new Error('Missing mandatory env variable API_URL');
}

function random(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1) + min);
}

async function* randomDelay(min: number, max: number): AsyncGenerator<number> {
    let interval = 0;
    while (true) {
        const ms = random(min, max);
        yield interval++;
        await delay(ms);
    }
}

function toQueryString<T extends Record<string, string | number>>(record: T): string {
    return Object.keys(record)
        .map((k) => `${k}=${record[k]}`)
        .join('&');
}

interface Request extends Record<string, number | string> {
    limit: number,
    str1: string,
    str2: string,
    int1: number,
    int2: number,
}

function fizzBuzzer(request: Request): Promise<Array<string>> {
    return fetch(`${apiUrl!}/api/v1/fizzbuzz?${toQueryString(request)}`)
        .then((response) => response.json().then((data) => [response.status, data] as const))
        .then(([status, data]) => {
            if (status !== 200) {
                throw new Error(data.error ?? `Server returned ${status}\n\t${data.message}`)
            } else {
                return data;
            }
        })
}

const strs = ['Fizz', 'Buzz', 'Lorem', 'Ipsum', ''];
const randomStr = () => {
    const idx = random(0, strs.length - 1);
    return strs[idx];
}

console.log('starting activity');

for await (const idx of randomDelay(1, 10)) {
    console.group(`request#${idx}`);
    try {
        const result = await fizzBuzzer({
            limit: random(-1, 10),
            str1: randomStr(),
            str2: randomStr(),
            int1: random(3, 5),
            int2: random(4, 6),
        });
        if (result !== undefined) {
            console.log(`got ${result.length} results back`);
        }
    } catch (err) {
        console.error(err);
    }
    console.groupEnd();
}

console.log('done');