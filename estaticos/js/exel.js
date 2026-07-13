// ===========================================================================
// Generador de Excel - Project EF
// ===========================================================================
// Quita el prefijo "data:image/png;base64," si existe

// En exel.js
const insertarGrafica = async (
  // <--- Asegúrate de que sea async
  workbook,
  sheet,
  base64Data,
  nombre,
  fila,
  columna,
  ancho,
  alto,
) => {
  if (!base64Data || typeof base64Data !== "string") {
    console.warn(`⚠️ La gráfica ${nombre} no tiene datos válidos.`);
    return;
  }

  try {
    const base64Limpio = base64Data.replace(/^data:image\/\w+;base64,/, "");

    // IMPORTANTE: workbook.addImage es síncrono pero a veces el buffer necesita tiempo
    const imageId = workbook.addImage({
      base64: base64Limpio,
      extension: "png",
    });

    sheet.addImage(imageId, {
      tl: { col: columna, row: fila },
      ext: { width: ancho, height: alto },
    });

    console.log(`✅ Gráfica ${nombre} insertada.`);
  } catch (err) {
    console.error(`❌ Error al insertar ${nombre}:`, err);
  }
};

const limpiarBase64 = (base64Data) => {
  return base64Data && base64Data.includes(",")
    ? base64Data.split(",")[1]
    : base64Data;
};

function limpiarMonto(valor) {
  if (valor === null || valor === undefined || valor === "" || valor === "-")
    return null;

  if (typeof valor === "number") return valor === 0 ? null : valor;

  let str = valor.toString().trim();
  let esNegativo = false;

  // Detectar formato contable (paréntesis)
  if (str.startsWith("(") && str.endsWith(")")) {
    esNegativo = true;
    str = str.substring(1, str.length - 1);
  }
  // Detectar signo menos estándar
  else if (str.startsWith("-")) {
    esNegativo = true;
  }

  // Limpiar comas y cualquier otro caracter no numérico
  let limpio = str.replace(/[^0-9.]/g, "");
  let num = parseFloat(limpio);

  if (isNaN(num) || num === 0) return null;

  return esNegativo ? num * -1 : num;
}

async function agregarPestañaConfig(workbook, nombreEmpresa) {
  const sheet = workbook.addWorksheet("Configuración");
  const estiloTitulo = {
    font: { bold: true, size: 14, color: { argb: "FFFFFF" } },
    fill: { type: "pattern", pattern: "solid", fgColor: { argb: "4A86E8" } },
    alignment: { horizontal: "center" },
  };
  const estiloEditable = {
    fill: { type: "pattern", pattern: "solid", fgColor: { argb: "FFF2CC" } },
    border: {
      top: { style: "thin" },
      left: { style: "thin" },
      bottom: { style: "thin" },
      right: { style: "thin" },
    },
  };

  sheet.mergeCells("B2:C2");
  sheet.getCell("B2").value = "CONFIGURACIÓN DEL REPORTE";
  sheet.getCell("B2").style = estiloTitulo;

  const datos = [
    ["Nombre de la Entidad:", nombreEmpresa || "NOMBRE DE TU EMPRESA"],
    ["Fecha de Inicio:", "01/01/2026"],
    ["Fecha de Cierre:", "31/12/2026"],
    ["Tasa ISR (%):", 0.25],
    ["Reserva Legal (%):", 0.05],
  ];

  datos.forEach((d, i) => {
    const row = i + 4;
    sheet.getCell(`B${row}`).value = d[0];
    sheet.getCell(`C${row}`).value = d[1];
    sheet.getCell(`C${row}`).style = estiloEditable;
    if (i >= 3) sheet.getCell(`C${row}`).numFmt = "0%";
  });

  sheet.getColumn("B").width = 30;
  sheet.getColumn("C").width = 50;
}

async function agregarPestañaCaratula(workbook, data, imageId) {
  const sheet = workbook.addWorksheet("Portada");

  // 1. Insertar Logo (Si existe)
  if (imageId !== null) {
    sheet.addImage(imageId, {
      tl: { col: 1, row: 1 }, // Columna B, Fila 2 (para dar un pequeño margen)
      ext: { width: 160, height: 160 },
    });
  }

  // 2. Configuración de textos
  // Nombre de la Empresa (Manual o desde data)
  const celdaNombre = sheet.getCell("B10");
  celdaNombre.value = data.nombreEmpresa;
  celdaNombre.font = { size: 18, bold: true };

  // Título del Reporte
  const celdaTitulo = sheet.getCell("B12");
  celdaTitulo.value = "ESTADOS FINANCIEROS";
  celdaTitulo.font = { size: 14, italic: true };

  // Referencia a la Configuración (Entidad)
  const celdaEntidad = sheet.getCell("B13");
  celdaEntidad.value = { formula: "'Configuración'!$C$4" }; // CORREGIDO: Usamos '='
  celdaEntidad.font = { size: 20, bold: true, color: { argb: "4A86E8" } };

  // 3. Estética: Anchar columna B para que quepa el nombre
  sheet.getColumn("B").width = 50;

  // Ocultar líneas de cuadrícula para que parezca una carátula real
  sheet.views = [{ showGridLines: false }];
}

async function agregarPestañaResultados(workbook, dataResultados, imageId) {
  const sheet = workbook.addWorksheet("Estado de Resultados");
  const azulEmpresarial = "4A86E8";

  // 1. INSERTAR LOGO (Solo una vez, tenías el código duplicado)
  if (imageId !== null) {
    sheet.addImage(imageId, {
      tl: { col: 0, row: 0 }, // Celda A1
      ext: { width: 60, height: 60 },
    });
  }

  // 2. TÍTULOS JUSTIFICADOS A LA DERECHA
  // Combinamos de A a D para tener todo el ancho del reporte
  sheet.mergeCells("A1:D1");
  const celdaEmpresa = sheet.getCell("A1");
  celdaEmpresa.value = { formula: "'Configuración'!$C$4" };
  celdaEmpresa.font = { bold: true, size: 14 };
  celdaEmpresa.alignment = { horizontal: "right", vertical: "middle" }; // <-- DERECHA

  sheet.mergeCells("A2:D2");
  const celdaFechas = sheet.getCell("A2");
  celdaFechas.value = {
    formula:
      '"ESTADO DE RESULTADOS DEL " & TEXT(\'Configuración\'!$C$5; "dd/mm/yyyy") & " AL " & TEXT(\'Configuración\'!$C$6; "dd/mm/yyyy")',
  };
  celdaFechas.font = { italic: true };
  celdaFechas.alignment = { horizontal: "right", vertical: "middle" }; // <-- DERECHA

  // 3. ENCABEZADOS DE TABLA
  const headers = ["DESCRIPCIÓN", "PARCIAL", "SUBTOTAL", "TOTAL"];
  const headerRow = sheet.getRow(4);
  headerRow.values = headers;
  headerRow.eachCell((c) => {
    c.font = { bold: true, color: { argb: "FFFFFF" } };
    c.fill = {
      type: "pattern",
      pattern: "solid",
      fgColor: { argb: azulEmpresarial },
    };
    c.alignment = { horizontal: "center" };
  });

  // 3. DATOS CON ESTILOS ESPECÍFICOS
  dataResultados.forEach((item, index) => {
    const rowIdx = index + 5;
    const row = sheet.getRow(rowIdx);
    const v = item.Value;

    if (!v || !v.nombre) return;

    // Llenado de datos
    row.getCell(1).value = v.nombre;
    row.getCell(2).value = limpiarMonto(v.col1);
    row.getCell(3).value = limpiarMonto(v.col2);
    row.getCell(4).value = limpiarMonto(v.col3);

    // --- FUNCIÓN INTERNA PARA APLICAR BORDES SOLO A CELDAS NUMÉRICAS ---
    const aplicarEstiloCelda = (celda, clase) => {
      if (!clase) return;

      let border = {};
      if (clase.includes("border-top")) border.top = { style: "thin" };
      if (clase.includes("border-bottom")) border.bottom = { style: "thin" };
      if (clase.includes("border-double")) border.bottom = { style: "double" };

      celda.border = border;

      // Si tiene línea o clase de negrita, ponemos el número en bold
      if (clase.includes("fw-bold") || clase.includes("border")) {
        celda.font = { bold: true };
      }
    };

    // --- APLICACIÓN DE ESTILOS SEGÚN LO QUE VIENE DE GO ---

    // 1. Estilo para el nombre (Columna 1)
    if (v.clasnombre && v.clasnombre.includes("titulo")) {
      row.getCell(1).font = { bold: true, color: { argb: "1F4E78" } };
    } else if (v.clasnombre && v.clasnombre.includes("fw-bold")) {
      row.getCell(1).font = { bold: true };
    }

    // 2. Estilos para los números (Columnas 2, 3 y 4)
    // Esto evita que las líneas se pasen a la descripción
    aplicarEstiloCelda(row.getCell(2), v.cla1); // Parcial
    aplicarEstiloCelda(row.getCell(3), v.cla2); // Subtotal
    aplicarEstiloCelda(row.getCell(4), v.cla3); // Total

    // 3. Resaltado especial si es el resultado final
    if (v.esresultado) {
      row.getCell(1).font = { bold: true, size: 12 };
    }

    // 4. FORMATO NUMÉRICO Y ALINEACIÓN
    const formatoContable = "#,##0.00;[Red](#,##0.00)";

    // Aplicarlo en tus bucles de celdas:
    for (let i = 2; i <= 5; i++) {
      const cell = row.getCell(i);
      cell.value = limpiarMonto(v["col" + (i - 1)]); // Asegúrate de pasar el valor limpio
      cell.numFmt = formatoContable;
      cell.alignment = { horizontal: "right" };
    }
  });

  sheet.getColumn(1).width = 45;
  sheet.getColumn(2).width = 15;
  sheet.getColumn(3).width = 15;
  sheet.getColumn(4).width = 15;
}

async function agregarPestañaHojaTrabajo(workbook, dataHoja) {
  const sheet = workbook.addWorksheet("Hoja de Trabajo");

  // Encabezados
  sheet.mergeCells("B1:C1");
  sheet.getCell("B1").value = "BALANCE DE SALDOS";
  sheet.mergeCells("D1:E1");
  sheet.getCell("D1").value = "ESTADO DE RESULTADOS";
  sheet.mergeCells("F1:G1");
  sheet.getCell("F1").value = "BALANCE GENERAL";

  const sub = [
    "DESCRIPCIÓN",
    "DEBE",
    "HABER",
    "PÉRDIDA",
    "GANANCIA",
    "ACTIVO",
    "PASIVO",
  ];
  sheet.getRow(2).values = sub;

  // 1. Definición de los estilos de relleno
  const fillSaldos = {
    type: "pattern",
    pattern: "solid",
    fgColor: { argb: "FFF5F5DC" },
  }; // Beige suave
  const fillResultados = {
    type: "pattern",
    pattern: "solid",
    fgColor: { argb: "FFE8F5E9" },
  }; // Verde pastel
  const fillBalance = {
    type: "pattern",
    pattern: "solid",
    fgColor: { argb: "FFE3F2FD" },
  }; // Celeste claro

  dataHoja.forEach((entrada) => {
    const datosFila = entrada.Value;
    const row = sheet.addRow([
      datosFila.nombre,
      limpiarMonto(datosFila.debe),
      limpiarMonto(datosFila.haber),
      limpiarMonto(datosFila.perdidas),
      limpiarMonto(datosFila.ganancias),
      limpiarMonto(datosFila.activo),
      limpiarMonto(datosFila.pasivo),
    ]);

    // 1. AHORA SÍ: Usamos la variable para identificar cualquier fila de cierre
    const nombreUpper = datosFila.nombre.toUpperCase();
    const esSumaFinal =
      nombreUpper.includes("SUMAS") || nombreUpper.includes("RESULTADO");

    for (let i = 2; i <= 7; i++) {
      const cell = row.getCell(i);
      const valor = cell.value;

      // Aplicar Colores (Beige, Verde, Celeste)
      if (valor !== null && valor !== 0 && valor !== "") {
        if (i <= 3) cell.fill = fillSaldos;
        else if (i <= 5) cell.fill = fillResultados;
        else if (i <= 7) cell.fill = fillBalance;
      }

      // 2. LÓGICA DE BORDES MEJORADA
      let estiloBorde = {
        top: { style: "thin", color: { argb: "E0E0E0" } },
        left: { style: "thin", color: { argb: "E0E0E0" } },
        bottom: { style: "thin", color: { argb: "E0E0E0" } },
        right: { style: "thin", color: { argb: "E0E0E0" } },
      };

      if (esSumaFinal) {
        cell.font = { bold: true };
        // Línea simple superior para todas las filas de totales
        estiloBorde.top = { style: "thin", color: { argb: "000000" } };

        // Si es específicamente la última fila, aplicamos la doble línea inferior
        if (nombreUpper.includes("SUMAS IGUALES")) {
          estiloBorde.bottom = { style: "double", color: { argb: "000000" } };
        }
      }

      cell.border = estiloBorde;
      cell.numFmt = "#,##0.00";
    }
  });
  sheet.getColumn(1).width = 40;
}

async function agregarPestañaBalance(workbook, dataBalance, imageId) {
  const sheet = workbook.addWorksheet("Balance General");

  // 1. INSERTAR LOGO
  if (imageId !== null) {
    sheet.addImage(imageId, {
      tl: { col: 0, row: 0 },
      ext: { width: 60, height: 60 },
    });
  }

  // 2. TÍTULOS (Asegúrate de que la hoja 'Configuración' exista)
  sheet.mergeCells("A1:E1");
  const celdaEmpresa = sheet.getCell("A1");
  celdaEmpresa.value = { formula: "'Configuración'!$C$4" };
  celdaEmpresa.font = { bold: true, size: 14 };
  celdaEmpresa.alignment = { horizontal: "right", vertical: "middle" };

  sheet.mergeCells("A2:E2");
  const celdaTitulo = sheet.getCell("A2");
  celdaTitulo.value = {
    formula:
      '"BALANCE GENERAL AL " & TEXT(\'Configuración\'!$C$6; "dd/mm/yyyy")',
  };
  celdaTitulo.font = { italic: true };
  celdaTitulo.alignment = { horizontal: "right", vertical: "middle" };

  // 4. LLENADO DE DATOS
  dataBalance.forEach((b) => {
    // Normalización de datos (Soporta mayúsculas y minúsculas de Go)
    const v = {
      nombre: b.Nombre || b.nombre || "",
      col1: b.Col1 !== undefined ? b.Col1 : b.col1,
      col2: b.Col2 !== undefined ? b.Col2 : b.col2,
      col3: b.Col3 !== undefined ? b.Col3 : b.col3,
      col4: b.Col4 !== undefined ? b.Col4 : b.col4,
      cla1: b.Cla1 || b.cla1 || "",
      cla2: b.Cla2 || b.cla2 || "",
      cla3: b.Cla3 || b.cla3 || "",
      cla4: b.Cla4 || b.cla4 || "",
      clasnombre: b.ClasNombre || b.clasnombre || "",
      esresultado: b.EsResultado || b.esresultado || false,
    };

    const row = sheet.addRow([
      v.nombre,
      limpiarMonto(v.col1),
      limpiarMonto(v.col2),
      limpiarMonto(v.col3),
      limpiarMonto(v.col4),
    ]);

    // --- FUNCIÓN PARA BORDES EN CELDAS NUMÉRICAS ---
    const aplicarEstiloCelda = (celda, clase) => {
      if (!clase) return;
      let border = {};
      if (clase.includes("border-top")) border.top = { style: "thin" };
      if (clase.includes("border-bottom")) border.bottom = { style: "thin" };
      if (clase.includes("border-double")) border.bottom = { style: "double" };
      celda.border = border;
      if (clase.includes("fw-bold") || clase.includes("border")) {
        celda.font = { bold: true };
      }
    };

    // --- APLICACIÓN DE ESTILOS ---

    // 1. Estilo para el nombre (Columna 1)
    if (v.clasnombre.includes("titulo")) {
      row.getCell(1).font = { bold: true, color: { argb: "000000" } }; // Negro sólido para Balance
    } else if (v.clasnombre.includes("fw-bold")) {
      row.getCell(1).font = { bold: true };
    }

    // 2. Estilos para los números (Columnas 2 a 5)
    aplicarEstiloCelda(row.getCell(2), v.cla1); // Parcial 1
    aplicarEstiloCelda(row.getCell(3), v.cla2); // Parcial 2
    aplicarEstiloCelda(row.getCell(4), v.cla3); // Total
    aplicarEstiloCelda(row.getCell(5), v.cla4); // Final

    // 3. Resaltado de sumas finales (Doble línea si es resultado)
    if (v.esresultado) {
      row.getCell(1).font = { bold: true };
      // El borde doble ya lo debería traer v.cla4 desde Go
    }

    // 4. Formatos de moneda y alineación
    for (let i = 2; i <= 5; i++) {
      const cell = row.getCell(i);
      cell.numFmt = "#,##0.00;[Red]-#,##0.00";
      cell.alignment = { horizontal: "right" };
    }
  });

  // 5. AJUSTES FINALES
  sheet.getColumn(1).width = 45;
  sheet.getColumn(2).width = 15;
  sheet.getColumn(3).width = 15;
  sheet.getColumn(4).width = 15;
  sheet.getColumn(5).width = 15;
  sheet.views = [{ showGridLines: false }];
}

async function agregarPestañaPartidasCierre(workbook, dataHoja) {
  const sheet = workbook.addWorksheet("Partidas de Cierre");

  // 1. Configuración de columnas y Encabezado Azul
  sheet.columns = [
    { header: "CUENTA / DESCRIPCIÓN", key: "desc", width: 45 },
    { header: "DEBE", key: "debe", width: 15 },
    { header: "HABER", key: "haber", width: 15 },
  ];

  const headerRow = sheet.getRow(1);
  headerRow.eachCell((cell) => {
    cell.font = { bold: true, color: { argb: "FFFFFF" } };
    cell.fill = {
      type: "pattern",
      pattern: "solid",
      fgColor: { argb: "2F5597" },
    };
    cell.alignment = { horizontal: "center" };
  });

  let sumaDebeP1 = 0,
    sumaHaberP1 = 0;
  let sumaDebeP2 = 0,
    sumaHaberP2 = 0;
  let invFinal = 0;

  // Estilo para las filas de "Pda."
  const estiloPda = (row) => {
    row.font = { italic: true, bold: true, color: { argb: "2F5597" } };
    row.getCell(1).border = {
      bottom: { style: "thin", color: { argb: "2F5597" } },
    };
  };

  // --- PDA 1: CIERRE DE RESULTADOS ---
  const pda1Header = sheet.addRow([
    "Pda. 1 - Liquidación de Cuentas de Resultado",
  ]);
  estiloPda(pda1Header);

  dataHoja.forEach((item) => {
    const v = item.Value;
    if (
      v.nombre.toUpperCase().includes("SUMA") ||
      v.nombre.toUpperCase().includes("RESULTADO")
    )
      return;

    if (v.nombre.toUpperCase() === "MERCADERIAS") {
      const invInicial = limpiarMonto(v.perdidas) || 0;
      invFinal = limpiarMonto(v.ganancias) || 0;
      if (invFinal > 0) {
        sheet.addRow([v.nombre + " (Inventario Final)", invFinal, null]);
        sumaDebeP1 += invFinal;
      }
      if (invInicial > 0) {
        sheet.addRow([
          "    " + v.nombre + " (Inventario Inicial)",
          null,
          invInicial,
        ]);
        sumaHaberP1 += invInicial;
      }
      return;
    }

    const perd = limpiarMonto(v.perdidas) || 0;
    const gan = limpiarMonto(v.ganancias) || 0;
    if (gan > 0) {
      sheet.addRow([v.nombre, gan, null]);
      sumaDebeP1 += gan;
    } else if (perd > 0) {
      sheet.addRow(["    " + v.nombre, null, perd]);
      sumaHaberP1 += perd;
    }
  });

  const utilidadEjec = sumaDebeP1 - sumaHaberP1;
  const filaRes = sheet.addRow([
    utilidadEjec > 0
      ? "    RESULTADO DEL EJERCICIO"
      : "RESULTADO DEL EJERCICIO",
    utilidadEjec < 0 ? Math.abs(utilidadEjec) : null,
    utilidadEjec > 0 ? utilidadEjec : null,
  ]);
  filaRes.font = { bold: true };

  if (utilidadEjec > 0) sumaHaberP1 += utilidadEjec;
  else sumaDebeP1 += Math.abs(utilidadEjec);

  // SUMAS IGUALES PDA 1
  const totalesP1 = sheet.addRow(["SUMAS IGUALES", sumaDebeP1, sumaHaberP1]);
  totalesP1.font = { bold: true };
  totalesP1.getCell(2).border = {
    top: { style: "thin" },
    bottom: { style: "double" },
  };
  totalesP1.getCell(3).border = {
    top: { style: "thin" },
    bottom: { style: "double" },
  };

  sheet.addRow([]); // Espacio

  // --- PDA 2: CIERRE DE BALANCE ---
  const pda2Header = sheet.addRow([
    "Pda. 2 - Cierre Final de Activos y Pasivos",
  ]);
  estiloPda(pda2Header);

  dataHoja.forEach((item) => {
    const v = item.Value;
    if (
      v.nombre.toUpperCase().includes("SUMA") ||
      v.nombre.toUpperCase().includes("RESULTADO")
    )
      return;

    if (v.nombre.toUpperCase() === "MERCADERIAS") {
      if (invFinal > 0) {
        sheet.addRow(["    " + v.nombre, null, invFinal]);
        sumaHaberP2 += invFinal;
      }
      return;
    }

    const act = limpiarMonto(v.activo) || 0;
    const pas = limpiarMonto(v.pasivo) || 0;
    if (pas > 0) {
      sheet.addRow([v.nombre, pas, null]);
      sumaDebeP2 += pas;
    } else if (act > 0) {
      sheet.addRow(["    " + v.nombre, null, act]);
      sumaHaberP2 += act;
    }
  });

  const filaRes2 = sheet.addRow([
    utilidadEjec > 0
      ? "RESULTADO DEL EJERCICIO"
      : "    RESULTADO DEL EJERCICIO",
    utilidadEjec > 0 ? utilidadEjec : null,
    utilidadEjec < 0 ? Math.abs(utilidadEjec) : null,
  ]);
  filaRes2.font = { bold: true };
  if (utilidadEjec > 0) sumaDebeP2 += utilidadEjec;
  else sumaHaberP2 += Math.abs(utilidadEjec);

  // SUMAS IGUALES PDA 2
  const totalesP2 = sheet.addRow(["SUMAS IGUALES", sumaDebeP2, sumaHaberP2]);
  totalesP2.font = { bold: true };
  totalesP2.getCell(2).border = {
    top: { style: "thin" },
    bottom: { style: "double" },
  };
  totalesP2.getCell(3).border = {
    top: { style: "thin" },
    bottom: { style: "double" },
  };

  // Formato numérico y alineación para todas las celdas de montos
  sheet.eachRow((row, rowNumber) => {
    if (rowNumber > 1) {
      row.getCell(2).numFmt = "#,##0.00";
      row.getCell(3).numFmt = "#,##0.00";
      row.getCell(2).alignment = { horizontal: "right" };
      row.getCell(3).alignment = { horizontal: "right" };
    }
  });
}

// En exel.js
async function agregarPestañaDashboard(
  workbook,
  dataIndices,
  logoId,
  graficas,
) {
  const sheet = workbook.addWorksheet("Dashboard Financiero");
  sheet.views = [{ showGridLines: false }];

  // 2. Título y Logo
  if (logoId) {
    sheet.addImage(logoId, {
      tl: { col: 0, row: 0 },
      ext: { width: 50, height: 50 },
    });
  }

  sheet.mergeCells("B2:E2");
  const t = sheet.getCell("B2");
  t.value = "ANÁLISIS DE INDICADORES Y PUNTO DE EQUILIBRIO";
  t.font = { bold: true, size: 14, color: { argb: "FFFFFF" } };
  t.fill = { type: "pattern", pattern: "solid", fgColor: { argb: "2F5597" } };
  t.alignment = { horizontal: "center" };

  // 3. Tabla de KPIs
  sheet.getRow(4).values = [
    "",
    "INDICADOR",
    "VALOR",
    "INTERPRETACIÓN",
    "DETALLE",
  ];
  sheet.getRow(4).font = { bold: true };

  let rowIdx = 5;
  dataIndices.forEach((kpi) => {
    const row = sheet.addRow([
      "",
      kpi.Nombre,
      kpi.Valor,
      kpi.Interpretacion,
      kpi.DetalleCuenta,
    ]);
    row.getCell(3).font = {
      bold: true,
      color: { argb: kpi.Valor.includes("-") ? "FF0000" : "0070C0" },
    };
    rowIdx++;
  });

  // Ajustar anchos
  sheet.getColumn(2).width = 30;
  sheet.getColumn(4).width = 45;
  sheet.getColumn(5).width = 45;

  let filaInsercion = 18;

  // AGREGAMOS UNA VALIDACIÓN EXTRA DE EMERGENCIA
  if (!window.imgEquilibrioBase64) {
    const can = document.getElementById("graficoEquilibrio");
    if (can) window.imgEquilibrioBase64 = can.toDataURL("image/png");
  }
  if (!window.imgDonaBase64) {
    const can = document.getElementById("graficoDona");
    if (can) window.imgDonaBase64 = can.toDataURL("image/png");
  }

  // 1. Gráfico de Equilibrio (Izquierda)
  if (window.imgEquilibrioBase64 && window.imgEquilibrioBase64.length > 100) {
    await insertarGrafica(
      workbook,
      sheet,
      window.imgEquilibrioBase64,
      "Equilibrio",
      filaInsercion, // Fila 18
      1, // Columna B (1)
      400, // Reducimos el ancho (antes 500)
      250, // Reducimos un poco el alto (antes 300)
    );
  }

  // 2. Gráfico de Dona (Derecha - A la par)
  if (window.imgDonaBase64 && window.imgDonaBase64.length > 100) {
    await insertarGrafica(
      workbook,
      sheet,
      window.imgDonaBase64,
      "Dona",
      filaInsercion, // Mismo nivel: Fila 18
      4, // Movemos a la Columna E (4) para que no se traslape
      250, // Reducimos el ancho para que sea cuadrada y pequeña
      250, // Reducimos el alto como pediste
    );
  }
}

// FUNCION MAESTRA
async function generarReporteExcelFinal() {
  const btnId = "btnExcel";
  const btn = document.getElementById(btnId);

  try {
    // 1. ESTADO DE CARGA
    if (btn) {
      btn.disabled = true;
      btn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Generando...';
    }

    // 2. CREAR EL LIBRO
    const wb = new ExcelJS.Workbook();

    // 3. REGISTRAR EL LOGO EN EL LIBRO
    let logoId = null;
    if (logoDataURL) {
      logoId = wb.addImage({
        base64: logoDataURL,
        extension: "png",
      });
    }

    // --- REGISTRAR LAS GRÁFICAS DEL DASHBOARD ---
    const canvasEquilibrio = document.getElementById("graficoEquilibrio");
    const canvasDona = document.getElementById("graficoDona");

    let graficasIds = { equilibrio: null, dona: null };

    if (canvasEquilibrio) {
      // LIMPIAMOS EL BASE64 ANTES DE AGREGARLO
      const base64Equilibrio = limpiarBase64(
        canvasEquilibrio.toDataURL("image/png"),
      );
      graficasIds.equilibrio = wb.addImage({
        base64: base64Equilibrio,
        extension: "png",
      });
    }

    if (canvasDona) {
      // LIMPIAMOS EL BASE64 ANTES DE AGREGARLO
      const base64Dona = limpiarBase64(canvasDona.toDataURL("image/png"));
      graficasIds.dona = wb.addImage({
        base64: base64Dona,
        extension: "png",
      });
    }
    // --------------------------------------------------

    // 4. GENERAR PESTAÑAS
    await agregarPestañaConfig(wb, dataReporte.nombreEmpresa);

    await agregarPestañaCaratula(
      wb,
      { nombreEmpresa: dataReporte.nombreEmpresa },
      logoId,
    );

    await agregarPestañaHojaTrabajo(wb, dataReporte.hojaTrabajo);

    await agregarPestañaResultados(wb, dataReporte.resultados, logoId);

    await agregarPestañaBalance(wb, dataReporte.balance, logoId);

    await agregarPestañaPartidasCierre(wb, dataReporte.hojaTrabajo);

    // Pasamos el objeto con los IDs de las gráficas
    await agregarPestañaDashboard(wb, dataReporte.indices, logoId, graficasIds);

    // 5. PROCESO DE DESCARGA
    const buffer = await wb.xlsx.writeBuffer();
    const fecha = new Date().toISOString().split("T")[0];
    const blob = new Blob([buffer], {
      type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    });

    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `Auditoria_${dataReporte.nombreEmpresa}_${fecha}.xlsx`;
    a.click();
    window.URL.revokeObjectURL(url);

    if (window.mostrarAviso) window.mostrarAviso("¡Excel generado con éxito!");
  } catch (error) {
    console.error("Error detallado:", error);
    alert("Fallo al generar Excel: " + error.message);
  } finally {
    // 6. RESTAURAR BOTÓN
    if (btn) {
      btn.disabled = false;
      btn.innerHTML =
        '<i class="fas fa-file-excel"></i> Generar Reporte de Auditoría (Excel)';
    }
  }
}
