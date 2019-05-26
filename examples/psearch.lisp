(load "search.lisp")

(defun ppresent (x lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        (t {or 
                (ppresent x (first-half lst))
                (ppresent x (second-half lst))} 
        )))

(defun smart-ppresent (x lst)
    (smart-ppresent-ric x lst (length lst)))

(defun smart-ppresent-ric (x lst initialSize)
    ; tries to balance the load
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        ((<= (length lst) (/ initialSize ncpu))
            (present x lst))
        (t {or 
                (smart-ppresent-ric x (first-half lst)  initialSize)
                (smart-ppresent-ric x (second-half lst) initialSize)
            } 
        )))

(defun genial-ppresent (x lst)
    (genial-ppresent-ric x lst 1))

(defun genial-ppresent-ric (x lst partitions)
    (cond 
        ((<= partitions ncpu)
            ; divide
            {or 
                (genial-ppresent-ric x (first-half lst)  (* partitions 2))
                (genial-ppresent-ric x (second-half lst) (* partitions 2))
            })
        (t (present x lst))
    ))

(write "[PAR] present first element...")
(time (ppresent 5900 llist))
(write "[PAR] present last element...")
(time (ppresent 9118 llist))
(write "")

(write "[PAR] smart present first element...")
(time (smart-ppresent 5900 llist))
(write "[PAR] smart present last element...")
(time (smart-ppresent 9118 llist))
(write "")

(write "[PAR] genial present first element...")
(time (genial-ppresent 5900 llist))
(write "[PAR] genial present last element...")
(time (genial-ppresent 9118 llist))
(write "")

t