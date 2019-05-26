(defun take (lst n) 
    (cond 
        ((eq n 0) nil) 
        (t (cons (car lst) (take (cdr lst) (- n 1))))))
        
(defun drop (lst n) 
    (cond ((eq n 0) lst) 
    (t (drop (cdr lst) (- n 1)))))
    
(defun merge (firstList secondList)
  (cond ((not firstList) secondList)
        ((not secondList) firstList)
        ((< (car firstList) (car secondList)) 
            (cons (car firstList) (merge (cdr firstList) secondList)))
        (t 
            (cons (car secondList) (merge firstList (cdr secondList))))))
            
(defun mergesort (lst)
  (cond 
    ((eq (length lst) 1) lst)
    (t (merge 
            (mergesort (take lst (/ (length lst) 2)))
            (mergesort (drop lst (/ (length lst) 2)))))))
            
(setq l '(4 2 56 73 3 1 94 3 1 23 32  2 34 32 1 2 4 52 23  5 4 765 87 6 35 13 321 12 432 3))

(write "Merge sorting...")
(time (dotimes (n 1000)(mergesort l)))
