(load "biglist.lisp")

(defun take (lst n) 
    (cond 
        ((eq n 0) nil) 
        (t (cons (car lst) (take (cdr lst) (- n 1))))))
        
(defun drop (lst n) 
    (cond ((eq n 0) lst) 
    (t (drop (cdr lst) (- n 1)))))

(defun first-half (lst)
    (take lst (/ (length lst) 2)))

(defun second-half (lst)
    (drop lst (/ (length lst) 2)))

(defun present (x lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        (t (or 
                (present x (first-half lst))
                (present x (second-half lst)))
        )))

(write "Sequential present")
(time (present 5900 llist))