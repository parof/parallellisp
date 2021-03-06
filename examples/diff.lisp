(defun first_arg (func) (car (cdr func)) )
(defun second_arg (func) (car (cdr (cdr func))) )
(defun fun_name (func) (car func) )
(defun mult (first second) (cons '* (cons first (cons second nil))) )
(defun expo (first second) (cons 'expt (cons first (cons second nil))) )
(defun sum  (first second) (cons '+ (cons first (cons second nil))) )
(defun diff (first second) (cons '- (cons first (cons second nil))) )
(defun frac (first second) (cons '/ (cons first (cons second nil))) )

(defun d (func) ( cond
    ((integerp func) 0) 
    ((symbolp func) 1)
    ((eq (fun_name func) 'ln) 
      (mult 
        (frac 1 (first_arg func)) 
        (d (first_arg func))))
    ((eq (fun_name func) 'exp)
      (mult 
        func 
        (d (first_arg func))))
    ((eq (car func) 'expt)
      (mult 
        (cons 'expt (cons (first_arg func) (- (second_arg func ) 1) ) ) 
        (second_arg func)))
    ((eq (fun_name func) 'sin)
      (mult
        (cons 'cos (cons (first_arg func) nil))
        (d (first_arg func))))
    ((eq (fun_name func) 'cos)
      (mult
        (mult (cons 'sin (cons (first_arg func) nil)) -1)
        (d (first_arg func))))
    ((eq (fun_name func) 'tan)
      (mult 
        (frac 
          1 
          (expo (cons 'cos (cons (first_arg func) nil)) 2) ) 
        (d (first_arg func))))
    ((eq (fun_name func) '+)
      (sum
        (d (first_arg func))
        (d (second_arg func))))
    ((eq (fun_name func) '*)
      (sum
        (mult (d (first_arg func)) (second_arg func))
        (mult (first_arg func) (d (second_arg func)))))
    ((eq (fun_name func) '/)
      (frac
        (diff 
          (mult (d (first_arg func)) (second_arg func)) 
          (mult (first_arg func) (d (second_arg func))))
        (expo (second_arg func) 2)))
    (t 
      nil)))

      ;;  (d '(/ (+ (expt x 2) 1) (cos x)))
      ;;  (d '(+ x 1))