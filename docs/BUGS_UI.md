# Bugs de UI y mejoras — `client_directory`

> Auditoría en vivo sobre `#client_directory` (45.973 pacientes en BD, 2026-11-19).

## Bugs encontrados

| # | Sev | Pantalla | Qué se ve en pantalla | Causa raíz verificada |
|---|---|---|---|---|
| CD-1 | 🔴 | `#client_directory` | **El campo `Edad` aparece en blanco al abrir un paciente que sí tiene fecha de nacimiento** (ej. `blanca rosa vivallo gomez`, `12/17/1945` → `Edad` vacío). | Bug compartido en `webtyp/components/agedate/agedate.go:150-198`: `ageInput` calcula `ageVal` una sola vez en `RenderInput` y nunca se suscribe al `SignalString` `value` cuando `form` carga un registro existente (F6 en el master). |
| CD-2 | 🟡 | `#client_directory` | **Placeholder mixto inglés/español `"Enter Dirección"`** en el campo `Dirección`. | `webtyp/input/address.go:21` llama `SetPlaceholder("Enter", "Address")` y `"Enter"` no está traducido (F3 en el master). |
| CD-3 | 🟡 | `#client_directory` | **Los chips de `Tipo de documento` y `Sexo` tienen un hueco vacío grande abajo** (`56px` con texto pegado arriba). | Bug compartido en `webtyp/form/css.go:105` (`PartRadioOption` — F1 en el master). |
| CD-4 | 🟡 | `#client_directory` | **Badges de fila truncados y con código técnico en minúscula** (`21514298-1 ·...`, `5739208-8 · rut`), y placeholder del buscador cortado (`Busca por documento, non`). | `ClientModel` concatena número + tipo en minúscula en el badge de `targetlist` + texto placeholder demasiado largo para el ancho del panel lateral. |
| CD-5 | ⚪ | `#client_directory` | **Al abrir la pantalla muestra `0 / 0` con lista gris vacía** (y botones de creación activos), sin un mensaje dentro de la lista que indique «Escribe nombre o RUT y presiona Enter para buscar entre los pacientes». | `ListsOnlyAfterSearch: true` (`client_directory/definition.go:30`) deja la lista en blanco al abrir y solo avisa «Presiona Enter para buscar» en un toast después de empezar a escribir. Además el input de fecha muestra formato `mm/dd/yyyy` dependiente del locale del navegador. |

## Sugerencias de mejora justificadas

1. **Mostrar estado guía dentro de la lista vacía en modo `ListsOnlyAfterSearch`:** en vez de una columna gris vacía con `0 / 0`, mostrar un texto centrado: «Busca por nombre, apellido o RUT (sin puntos) y presiona Enter». *Justificación:* un recepcionista nuevo no sabe si hay 0 pacientes o si debe buscar primero.
2. **Formatear el badge de documento con etiqueta corta en mayúscula (`RUT`, `PAS`, `CÉD`) o mostrar el número limpio sin saturar el ancho:** *Justificación:* hoy los RUT de 8 dígitos + dígito verificador se cortan justo en `· ...`, ocultando el tipo de documento en casi todas las filas.
3. **Acortar el placeholder del buscador a «Nombre, apellido o RUT…»:** *Justificación:* cabe completo en desktop (`320px`) y móvil (`390px`) sin truncarse en `non`.
