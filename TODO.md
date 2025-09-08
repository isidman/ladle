# TODO List for "Ladle" project

### 1. Back-end building:
 - Create structs for all color formats showing in ladle.

 - ColorRequest struct must be made to represent the colors requested from the web app client (end user).

 - PaletteRequest struct must be made to represent the type of palette the end user would like to create based on the color they chose from the spectrum.

 - Conversion control for color formats.

 - Opening the web app makes a random color generation.
 
 - Create the color conversion algorithm for each type.

 
 - Specific color conversion & Palette generation functions should be added at the end of the backend file for easier reference.


### 2. JS Application: 
  - Create pre-initialization state and get the element ids from backend.

  - Initialization state with color conversion utility setup.

  - Setup canvas drawing and & clear.

  - Draw in the colorwheel.

  - Setup event handling, inputs, sliders and buttons.

  - Handle updates in brightness slider and  input fields. *Check documentation on how to do this without triggering events.*

  - Storage management

  - Clipboard functionality and fallback handling for browsers that can't support it.

 - Check DOM loading and run diagnostics.
