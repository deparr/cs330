const t1 = {
	l: {
		l: {
			l: {
				v: 1
			},
			r: {
				v: 2
			}
		},
		r: {
			l: {
				v: 3
			},
			r: {
				v: 4
			}
		}
	},
	r: {
		l: {
			v: 5
		},
		r: null
	}
};

const t3 = {
	l: {
		l: {
			l: {
				v: 1
			},
			r: {
				v: 2
			}
		},
		r: {
			l: {
				v: 4
			},
			r: {
				v: 4
			}
		}
	},
	r: {
		l: null,
		r: {
			v: 5
		}
	}
};

function chan() {
	let val = undefined;
	return {
		in: async (new_val) => {
			while (val != undefined) {
				await new Promise(r => setTimeout(r, 10));
			}
			val = new_val;
		},
		out: async () => {
			while (val == undefined) {
				await new Promise(r => setTimeout(r, 10));
			}
			const ret = val;
			val = undefined;
			return ret;
		},
	}
}

export function coroutine(f) {
	const cin = chan();
	const cout = chan();
	const resume = async (inv) => {
		await cin.in(inv);
		return await cout.out();
	}

	const _yield = async (outv) => {
		await cout.in(outv);
		return await cin.out();
	};

	(async () => {await cout.in(f(await cin.out(), _yield))})()
	return resume
}

// let tree = 1;
// function coroutine() {
// 	let next_value = undefined;
// 	let tree_inner = tree++;
//
// 	const set = async (val) => {
// 		console.log(`setting ${tree_inner} ${val}`);
// 		while (next_value != undefined) {
// 			console.log(`waiting set ${tree_inner} ${val}, ${next_value}`)
// 			await new Promise(r => setTimeout(r, 50));
// 		}
// 		console.log(`set tree ${tree_inner} ${val}`);
// 		next_value = val;
// 	}
//
// 	const get = async () => {
// 		while (next_value == undefined) {
// 			console.log(`waiting get ${tree_inner}, ${next_value}`)
// 			await new Promise(r => setTimeout(r, 50));
// 		}
// 		const ret = next_value;
// 		next_value = undefined;
// 		return ret;
// 	}
//
// 	return [set, get];
// }

// async function visit2(t, _yield) {
// 	if (t) {
// 		if (t.l) {
// 			await visit2(t.l, _yield);
// 		} else if (t.r) {
// 			await visit2(t.r, _yield);
// 		} else if (t.v) {
// 			await _yield(t.v);
// 		}
// 	}
// }
//
// async function visit(t, _yield) {
// 	if (t.l || t.r) {
// 		if (t.l)
// 		await visit2(t.l, _yield)
// 		if (t.r)
// 		await visit2(t.r, _yield)
// 	} else if (t.v) {
// 		await _yield(t.v)
// 	}
// 	await _yield(null);
// 	return;
// }
//

// async function cmp(t1, t2) {
// 	let [setCo1, getCo1] = coroutine(visit, t1)
// 	let [setCo2, getCo2] = coroutine(visit, t2)
// 	let iter = 1;
//
// 	visit(t1, setCo1);
// 	visit(t2, setCo2);
// 	while (true & iter < 200) {
// 		console.log('waiting on co1...')
// 		let v1 = await getCo1()
// 		console.log('waiting on co2...')
// 		let v2 = await getCo2()
//
// 		console.log("got vals", v1, v2);
// 		if (v1 == undefined != v2 == undefined) {
// 			return false
// 		}
//
// 		if (v1 != v2) {
// 			return false
// 		}
//
// 		if (v1 == null && v2 == null) {
// 			return true
// 		}
// 		iter++;
// 	}
// }

// const res = await cmp(t1, t3);

// export default res;

