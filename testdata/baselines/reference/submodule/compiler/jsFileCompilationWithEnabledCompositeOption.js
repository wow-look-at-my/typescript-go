//// [tests/cases/compiler/jsFileCompilationWithEnabledCompositeOption.ts] ////

//// [a.ts]
class c {
}

//// [b.js]
function foo() {
}


//// [out.js]
"use strict";
class c {
}
function foo() {
}
