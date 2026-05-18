// ==========================================
// 1. VARIABLES Y SUMADORA
// ==========================================
let currentInput = '0';
let subtotal = 0;
let pendingOp = null;
let ultimoInputEnfocado = null;
let yaSonó = false;
const sonidoExito = new Audio('https://assets.mixkit.co/active_storage/sfx/2568/2568-preview.mp3');
sonidoExito.volume = 0.3;

window.updateDisplay = function() {
  const display = document.getElementById('displayCalc');
  if (display) display.value = currentInput;
};

window.addNum = function(num) {
  if (currentInput === '0') currentInput = num;
  else currentInput += num;
  window.updateDisplay();
};

window.addDecimal = function() {
  if (!currentInput.includes('.')) {
    currentInput += (currentInput === '' ? '0.' : '.');
    window.updateDisplay();
  }
};

window.backspace = function() {
  if (currentInput.length > 1) {
    currentInput = currentInput.slice(0, -1);
  } else {
    currentInput = '0';
  }
  window.updateDisplay();
};

window.clearCalc = function() {
  currentInput = '0';
  subtotal = 0;
  pendingOp = null;
  const cinta = document.getElementById('cintaContable');
  if (cinta) cinta.innerHTML = '<div class="text-muted small text-center border-bottom mb-1">REGISTRO</div>';
  window.updateDisplay();
};

window.setOp = function(op) {
  const valorActual = parseFloat(currentInput);
  if (pendingOp !== null) window.ejecutarOperacion(valorActual);
  else subtotal = valorActual;
  pendingOp = op;
  window.addToCinta(currentInput + ' ' + op);
  currentInput = '0';
  window.updateDisplay();
};

window.ejecutarOperacion = function(nuevoValor) {
  if (pendingOp === '+') subtotal += nuevoValor;
  else if (pendingOp === '-') subtotal -= nuevoValor;
  else if (pendingOp === '*') subtotal *= nuevoValor;
  else if (pendingOp === '/') subtotal = nuevoValor !== 0 ? subtotal / nuevoValor : 0;
};

window.calcular = function() {
  if (pendingOp === null) return;
  const valorFinal = parseFloat(currentInput);
  window.ejecutarOperacion(valorFinal);
  window.addToCinta(currentInput);
  window.addToCinta("----------");
  window.addToCinta("* " + subtotal.toLocaleString('en-US', { minimumFractionDigits: 2 }));
  currentInput = subtotal.toString();
  pendingOp = null;
  window.updateDisplay();
};

window.toggleCalc = function() {
  const calc = document.getElementById('calculadoraContable');
  const btnFlotante = document.getElementById('btnSumadora');
  if (!calc) return;

  if (calc.classList.contains('d-none')) {
    calc.classList.remove('d-none');
    if (btnFlotante) btnFlotante.style.display = 'none';
    document.querySelectorAll('.input-contable').forEach(i => {
      if (i.id !== 'input-220001') i.style.backgroundColor = "#1e1e1e";
    });
  } else {
    calc.classList.add('d-none');
    if (btnFlotante) btnFlotante.style.display = 'block';
  }
};

window.pegarResultado = function() {
  const display = document.getElementById('displayCalc');
  if (display && ultimoInputEnfocado) {
    ultimoInputEnfocado.value = display.value;
    ultimoInputEnfocado.dispatchEvent(new Event('input'));
    window.addToCinta("-> Trasladado");
    ultimoInputEnfocado.focus();
  } else {
    window.mostrarAviso("Primero selecciona una celda");
  }
};

window.addToCinta = function(text) {
  const cinta = document.getElementById('cintaContable');
  if (!cinta) return;
  const line = document.createElement('div');
  line.style.color = text.includes('-') ? "#ff4444" : "inherit";
  line.style.fontSize = "0.75rem";
  line.innerText = text;
  cinta.appendChild(line);
  cinta.scrollTop = cinta.scrollHeight;
};

window.mostrarAviso = function(mensaje) {
  const toastEl = document.getElementById('toastAtajo');
  if (!toastEl) { console.log(mensaje); return; }
  document.getElementById('toastMensaje').innerHTML = `<i class="fas fa-keyboard me-2"></i> ${mensaje}`;
  new bootstrap.Toast(toastEl, { delay: 2000 }).show();
};

// ==========================================
// 2. LÓGICA DE CUADRE Y EVENTOS DE TABLA
// ==========================================
const formatearVisual = (input) => {
  let val = parseFloat(input.value.replace(/,/g, ''));
  if (!isNaN(val)) {
    input.value = val.toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    });
  }
};

const calcularDiferencia = () => {
  let totalDebe = 0;
  let totalHaber = 0;
  let llenados = 0;
  const inputs = document.querySelectorAll('.input-contable');
  const diffValue = document.getElementById('diffValue');
  const cuadro = document.getElementById('cuadroCuadre');

  inputs.forEach(input => {
    const nameInput = input.getAttribute('name') || "";
    const valorLimpio = input.value.replace(/,/g, '');
    const val = parseFloat(valorLimpio) || 0;

    if (nameInput === "220001" || nameInput === "220005") return;

    if (val !== 0) {
      llenados++;
      const saldo = input.getAttribute('data-saldo');
      if (['activo', 'perdida', 'costo_debe'].includes(saldo)) totalDebe += val;
      else if (['pasivo', 'ganancia', 'costo_haber'].includes(saldo)) totalHaber += val;
    }
  });

  const diff = Math.abs(totalDebe - totalHaber);
  if (diffValue && cuadro) {
    diffValue.innerText = diff.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
    if (diff < 0.01 && llenados > 0) {
      diffValue.style.color = "#00ff00";
      diffValue.innerHTML = '<i class="fas fa-check-circle"></i> CUADRADO';
      if (!yaSonó) { sonidoExito.play().catch(() => { }); yaSonó = true; }
    } else {
      diffValue.style.color = "#ff4444";
      yaSonó = false;
    }
  }
};

document.addEventListener('DOMContentLoaded', () => {
  const cuadro = document.getElementById('cuadroCuadre');
  if (cuadro) {
    document.body.appendChild(cuadro);
    cuadro.style.setProperty("right", "10px", "important");
    cuadro.style.setProperty("left", "auto", "important");
  }

  document.addEventListener('input', (e) => {
    if (e.target.classList.contains('input-contable')) {
      calcularDiferencia();
      const mappings = { '111301': '220001', '111303': '311001', '111304': '311002', '111305': '311001' };
      if (mappings[e.target.name]) {
        const target = document.getElementsByName(mappings[e.target.name])[0];
        if (target) { target.value = e.target.value; target.dispatchEvent(new Event('input')); }
      }
    }
  });

  document.addEventListener('focusin', (e) => {
    if (e.target.classList.contains('input-contable')) {
      ultimoInputEnfocado = e.target;
      e.target.value = e.target.value.replace(/,/g, '');
      e.target.select();
    }
  });

  document.addEventListener('focusout', (e) => {
    if (e.target.classList.contains('input-contable')) formatearVisual(e.target);
  });

  calcularDiferencia();
  // Interceptar el clic en el botón de submit para limpiar comas
  const form = document.getElementById('balanceForm');
  if (form) {
    form.addEventListener('submit', function(e) {
      console.log("🚀 Limpiando datos para envío...");
      document.querySelectorAll('.input-contable').forEach(input => {
        // Quitamos comas antes de que viajen al servidor
        input.value = input.value.replace(/,/g, '');
      });
    });
  }
});


// ==========================================
// 3. ATAJOS DE TECLADO (Ctrl+S, Alt+C, etc)
// ==========================================
document.addEventListener('keydown', (e) => {
  const isS = (e.key === 's' || e.key === 'S' || e.code === 'KeiS');
  if (!e || !e.key) return;
  const key = e.key.toLowerCase();
  const isCtrl = e.ctrlKey || e.metaKey;
  const enInputTabla = e.target.classList.contains('input-contable');
  const calc = document.getElementById('calculadoraContable');
  const estaVisible = calc && !calc.classList.contains('d-none');

  if (isCtrl && isS) {
    e.preventDefault(); // DETIENE el "Guardar como" del navegador
    e.stopPropagation(); // EVITA que otros scripts lo vean
    console.log("💾 Iniciando guardado contable...");
    // 1. Limpiar comas de todos los inputs antes de enviar
    document.querySelectorAll('.input-contable').forEach(input => {
      input.value = input.value.replace(/,/g, '');
    });
    // 2. Localizar el formulario
    const form = document.getElementById('balanceForm');
    const diffValue = document.getElementById('diffValue');
    const estaCuadrado = diffValue && diffValue.innerText.includes("CUADRADO");

    if (form) {
      if (estaCuadrado) {
        window.mostrarAviso("¡Cuadrado! Guardando datos...");
        form.submit();
      } else {
        if (confirm("⚠️ La hoja NO está cuadrada. ¿Deseas guardar de todas formas?")) {
          form.submit();
        }
      }
    } else {
      console.error("No se encontró el formulario #balanceForm");
    }
    return false;
  }

  // Enter en tabla
  if (key === 'enter' && enInputTabla) {
    e.preventDefault();
    const inputs = Array.from(document.querySelectorAll('.input-contable'));
    const index = inputs.indexOf(e.target);
    if (index > -1 && index < inputs.length - 1) inputs[index + 1].focus();
    return;
  }

  // Alt + C
  if (e.altKey && key === 'c') {
    e.preventDefault();
    window.toggleCalc();
    return;
  }

  // Teclas calculadora
  if (estaVisible && !enInputTabla) {
    if (/[0-9]/.test(key)) {
      e.preventDefault();
      window.addNum(key);
    }

    if (['+', '-', '*', '/'].includes(key)) {
      e.preventDefault();
      window.setOp(key);
    }

    if (key === '.' || key === ',') {
      e.preventDefault();
      window.addDecimal();
    }

    if (key === 'backspace') {
      e.preventDefault();
      window.backspace();
    }
    if (key === 'enter') {
      e.preventDefault();
      window.calcular();
    }

    if (key === 'p') {
      e.preventDefault();
      window.pegarResultado();
    }

    if (key === 'c' && !e.altKey) {
      e.preventDefault();
      window.clearCalc();
    }

    if (key === 'escape') {
      e.preventDefault();
      window.toggleCalc();
    }

  }
}, { capture: true });

