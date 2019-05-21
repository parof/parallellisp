(defun fib (n) (cond ((eq n 0) 0) ((eq n 1) 1) (t (+ (fib (- n 1)) (fib (- n 2))))) )
(defun p-fib (n) (cond ((eq n 0) 0) ((eq n 1) 1) (t {+ (fib (- n 1)) (fib (- n 2))})) )
(defun bench (fun n) (+ (fun n) (fun n) (fun n) (fun n) (fun n) (fun n) (fun n) (fun n)))
(defun p-bench (fun n) {+ (fun n) (fun n) (fun n) (fun n) (fun n) (fun n) (fun n) (fun n)})
(time (bench fib 25))
(time (bench p-fib 25))
(time (p-bench fib 25))
(time (p-bench p-fib 25))