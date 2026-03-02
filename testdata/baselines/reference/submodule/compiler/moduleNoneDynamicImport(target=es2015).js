//// [tests/cases/compiler/moduleNoneDynamicImport.ts] ////

//// [a.ts]
const foo = import("./b");

//// [b.js]
export default 1;


//// [a.js]
export default 1;
const foo = import("./b");
