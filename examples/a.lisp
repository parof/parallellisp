((lambda (x y) (x y)) 'car '(1 2 3 4e))
((lambda (x y) (+ x y)) 1 2)
(time (+ 1 2 3 4 5 5 6 6 74345 35 34 5 45 245 423545 4 523 454 2352 ))

(defun inc (x) (cond ((eq x 900) t) (t (inc (+ x 1)))))