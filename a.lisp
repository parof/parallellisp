((lambda (x y) (x y)) 'car '(1 2 3 4e))
((lambda (x y) (+ x y)) 1 2)
(time (+ 1 2 3 4 5 5 6 6 74345 35 34 5 45 245 423545 4 523 454 2352 ))

(defun inc (x) (cond ((eq x 900) t) (t (inc (+ x 1)))))


(defun bench (n) {+ (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) )
(time {+ (fib 25) (fib 25) (fib 25) (fib 25) (fib 25) (fib 25) (fib 25) (fib 25) } )
(time (+ (fib 25) (fib 25) (fib 25) (fib 25) (fib 25) (fib 25) (fib 25) (fib 25) ) )

(defun toz (n) (cond ((eq n 0) 0) (t (toz (- n 1))) ) )

(defun inc (n) (cond ((eq n 1000) n) (t (inc (+ n 1))) ))
(defun fib (n) (cond ((eq n 0) 0) ((eq n 1) 1) (t (+ (fib (- n 1)) (fib (- n 2))))) )
