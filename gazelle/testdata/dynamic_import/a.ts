
let lamdba = () => import('./b')
lamdba()

// some comment
/*
should be ignored:
import('./d')
*/

let lamdba2 = () => import('./c')
lamdba2()

