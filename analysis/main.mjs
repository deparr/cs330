import {coroutine} from "./control.mjs";

const t1 = { l: { l: { l: { v: 1 }, r: { v: 2 } }, r: { l: { v: 3 }, r: { v: 4 } } }, r: { l: { v: 5 }, r: null } };
const t3 = { l: { l: { l: { v: 1 }, r: { v: 2 } }, r: { l: { v: 4 }, r: { v: 4 } } }, r: { l: null, r: { v: 5 } } };

async function cmp() {
}

const res = coroutine();

