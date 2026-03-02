//// [tests/cases/compiler/jsFileCompilationDuplicateVariableErrorReported.ts] ////

//// [b.js]
var x = "hello";

//// [a.ts]
var x = 10; // Error reported


//// [out.js]
"use strict";
var x = "hello";
var x = 10; // Error reported
