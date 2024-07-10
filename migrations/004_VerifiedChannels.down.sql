DELETE FROM clients_verified
WHERE (
    hash = 'b4ec3c4334a0249dae95c284ec5983df' AND type = 0 OR
    hash = '74be16979710d4c4e7c6647856088456' AND type = 0 OR
    hash = 'd41d8cd98f00b204e9800998ecf8427e' AND type = 0 OR
    hash = 'ad921d60486366258809553a3db49a4a' AND type = 1 OR
    hash = '74be16979710d4c4e7c6647856088456' AND type = 1 OR
    hash = 'd41d8cd98f00b204e9800998ecf8427e' AND type = 1 OR
    hash = 'ad921d60486366258809553a3db49a4a' AND type = 2 OR
    hash = 'dcfcd07e645d245babe887e5e2daa016' AND type = 2 OR
    hash = '28c8edde3d61a0411511d3b1866f0636' AND type = 2 OR
    hash = '74be16979710d4c4e7c6647856088456' AND type = 2 OR
    hash = 'd41d8cd98f00b204e9800998ecf8427e' AND type = 2 OR
    hash = 'd1c651c36f499849f1c9a5843567e686' AND type = 2
);