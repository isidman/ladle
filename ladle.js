class ColorPicker {
//PRE-INITIALIZATION
    constructor() {
        // App data
        this.wheelSize = 280;
        this.maxRecentColors = 8;
        this.maxSavedColors = 10;

        //Current color state
        this.currentColor = {
            hsv: {h: 0, s: 100, v: 100 },
            rgb: {r: 255, g: 0, b: 0 },
            hsl: {h: 0, s: 100, l: 50 },
            hex: '#ff0000'
        };

        //State for input handling
        this.isUpdatingFromUI = false;

        //DOM elements
        this.canvas = document.getElementById('colorWheel');
        this.ctx = this.canvas.getContext('2d');
        this.indicator = document.getElementById('wheelIndicator');
        this.brightnessSlider = document.getElementById('brightnessSlider');
        this.colorPreview = document.getElementById('colorPreview');

        //Input elements
        this.inputs = {
            hex: document.getElementById('hexInput'),
            
            rgb: {
                r: document.getElementById('rgbR'),
                g: document.getElementById('rgbG'),
                b: document.getElementById('rgbB')
            },

            hsl: {
                h: document.getElementById('hslH'),
                s: document.getElementById('hslS'),
                l: document.getElementById('hslL')
            },

            hsv: {
                h: document.getElementById('hsvH'),
                s: document.getElementById('hsvS'),
                v: document.getElementById('hsvV')
            }
        };

        //Storage containers
        this.recentColorsContainer = 
        document.getElementById('recentColors');
        this.savedPaletteContainer =
        document.getElementById('savedPalette');

        //INITIALIZATION
        this.init();
    }

}