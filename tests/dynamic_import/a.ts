
let lamdba = () => import('./b').do_something("something")
lamdba()

// some comment
/*
should be ignored:
import('./d')
*/

let lamdba2 = () => import('./c')
lamdba2()

