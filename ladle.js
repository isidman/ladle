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

    init() {
        this.drawColorWheel();
        this.setupEventListeners();
        this.loadStoredColors();
        this.updateUI();
        this.updateIndicatorPosition();
    }

    //Color conversion utilities
    hsvToRgb(h, s, v) {
        h = h / 360;
        s = s / 100;
        v = v / 100;

        const c = v * s;
        const x = c *(1 - Math.abs((h * 6) % 2 - 1));
        const m = v - c;

        let r, g, b;

        if (h >= 0 && h < 1/6) {
            r=c; g=x; b=0;
        } else if ( h >= 1/6 && h < 2/6) {
            r=x; g=c; b=0;
        } else if ( h >= 2/6 && h < 3/6) {
            r=0; g=c; b=x;
        } else if ( h >= 3/6 && h < 4/6) {
            r=0; g=x; b=c;
        } else if ( h >= 4/6 && h < 5/6) {
            r=x; g=0; b=c;
        } else {
            r=c; g=0; b=x;
        }

        return {
            r: Math.round((r + m) * 255),
            g: Math.round((g + m) * 255),
            b: Math.round((b + m) * 255)
        };
    }

    rgbToHsv(r, g, b) {
        r /= 255;
        g /= 255;
        b /= 255;

        const max = Math.max(r, g, b);
        const min = Math.min(r, g, b);
        const diff = max - min;

        let h = 0;
        const s = max === 0 ? 0 : diff / max;
        const v = max;

        if (diff !== 0) {
            if (max === r) {
                h = (60 * ((g - b) / diff ) + 360) % 360;
            } else if (max === g) {
                h = (60 * ((b -r) / diff) + 120) % 360;
            }else {
                h = (60 * ((r - g) / diff) + 240) % 360;
            }
        }

        return {
            h: Math.round(h),
            s: Math.round(s * 100),
            v: Math.round(v * 100)
        };
    }
    
    rgbToHsl(r, g, b) {
        r /= 255;
        g /= 255;
        b /= 255;

        const max = Math.max(r, g, b);
        const min = Math.min(r, g, b);
        const diff = max - min;

        const l = (max + min) / 2;
        let h = 0;
        let s = 0;

        if (diff !== 0) {
            s = l > 0.5 ? diff / (2 - max - min) : diff / (max + min);
                if (max === r) {
                    h = (60 * ((g - b) / diff) + 360) % 360;
                } else if (max === g) {
                    h = (60 * ((b - r) / diff) + 120) % 360;
                } else {
                    h = (60 * ((r - g) / diff) + 240) % 360;
                }
        }

        return {
            h: Math.round(h),
            s: Math.round(s* 100),
            l: Math.round(l * 100)
        };
    }

    hslToRgb(h, s, l) {
        h /= 360;
        s /= 100;
        l /= 100;

        const c = (1 - Math.abs(2 * l -1)) * s;
        const x = c * (1 - Math.abs((h * 6) % 2 - 1));
        const m = l - c / 2;

        let r, g, b;

        if (h >= 0 && h < 1/6) {
            r=c; g=x; b=0;
        } else if (h >= 1/6 && h < 2/6) {
            r=x; g=c; b=0;
        } else if (h >= 2/6 && h < 3/6) {
            r=0; g=c; b=x;
        } else if (h >= 3/6 && h < 4/6) {
            r=0; g=x; b=c;
        } else if (h >= 4/6 && h < 5/6) {
            r=x; g=0; b=c;
        } else {
            r=c; g=0; b=x;
        }

        return {
            r: Math.round((r + m) * 255),
            g: Math.round((g + m) * 255),
            b: Math.round((b + m) * 255)
        };
    }

    rgbToHex(r, g, b) {
        return '#' + [r, g, b].map(x => {
           const  hex = x.toString(16);
           return hex.length === 1 ? '0' + hex : hex; 
        }).join('');
    }

    hexToRgb(hex) {
        let r = 0, g = 0, b = 0;
        const offset = hex.startsWith('#') ? 1 : 0;

        if (hex.length === 4 + offset) {
            //3-digit hex
            r = parseInt(hex[offset] + hex[offset], 16);
            g = parseInt(hex[offset + 1] + hex[offset + 1], 16);
            b = parseInt(hex[offset + 2] + hex[offset + 2], 16);
        } else if (hex.length === 7 + offset) {
            //6-digit hex
            r = parseInt(hex[offset] + hex[offset + 1], 16);
            g = parseInt(hex[offset + 2] + hex[offset + 3], 16);
            b = parseInt(hex[offset + 4] + hex[offset + 5], 16);
        }

        return {r, g, b};
    }

    //Canvas drawing
    drawColorWheel() {
        const centerX = this.wheelSize / 2;
        const centerY = this.wheelSize / 2;
        const radius = this.wheelSize / 2 - 10;

        //Clear canvas
        this.ctx.clearRect(0, 0, this.wheelSize, this.wheelSize);

        //Draw the color wheel
        const imageData = 
        this.ctx.createImageData(this.wheelSize, this.wheelSize);
        const data = imageData.data;

        for (let x = 0; x < this.wheelSize; x++) {
            for (let y = 0; y < this.wheelSize; y++) {
                const dx = x - centerX;
                const dy = y - centerY;
                const distance = Math.sqrt(dx * dx + dy * dy);

                if (distance <= radius) {
                    const angle = Math.atan2(dy, dx) * 180 / Math.PI;
                    const hue = (angle + 360) % 360;
                    const saturation = (distance / radius) * 100;
                    const value = 100;
                    const rgb = this.hsvToRgb(hue, saturation, value);
                    const index = (y * this.wheelSize + x) * 4;
                    data[index] = rgb.r;
                    data[index + 1] = rgb.g;
                    data[index + 2] = rgb.b;
                    data[index + 3] = 255;
                } else {
                    const index = (y * this.wheelSize + x) * 4;
                    data[index + 3] = 0; //Transparency
                }
            }
        }

        this.ctx.putImageData(imageData, 0, 0);
    }

    //Event handling
    setupEventListeners() {
        //Click and drag on the canvas
        this.canvas.addEventListener('mousedown', (e) => {
            this.handleWheelInteraction(e);
            this.canvas.addEventListener('mousemove', this.handleMouseMove);
        });

        this.canvas.addEventListener('mouseup', () =>{
            this.canvas.removeEventListener('mousemove', this.handleMouseMove);
        });

        //Store bound function for removal
        this.handleMouseMove = (e) =>
            this.handleWheelInteraction(e);

        //Brightness slider
        this.brightnessSlider.addEventListener('input', (e) => {
            this.currentColor.hsv.v = parseInt(e.target.value);
            this.updateFromHSV();
        });

        //Input field changes
        this.inputs.hex.addEventListener('input', () => this.updateFromHex());
        this.inputs.hex.addEventListener('change', () => this.updateFromHex());

        Object.values(this.inputs.rgb).forEach(input => {
            input.addEventListener('input', () => this.updateFromRGB());
            input.addEventListener('change', () => this.updateFromRGB());
        });

        Object.values(this.inputs.hsl).forEach(input => {
            input.addEventListener('input', () => this.updateFromHSL());
            input.addEventListener('change', () => this.updateFromHSL());
        });

        Object.vaules(this.inputs.hsv). forEach(input => {
            input.addEventListener('input', () => this.updateFromHSV());
            input.addEventListener('change', () => this.updateFromHSV());
        });

        //Copy buttons
        document.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.copyToClipboard(e.target.dataset.copy));
        });

        //Save and clear buttons
        document.getElementById('saveColorBtn').addEventListener('click', () => this.saveCurrentColor());
        document.getElementById('clearRecentBtn').addEventListener('click', () => this.clearRecentColors());
    }
    
    handleWheelInteraction(e) {
        const rect = this.canvas.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        

        const centerX = this.wheelSize / 2;
        const centerY = this.wheelSize / 2;
        const dx = x - centerX;
        const dy = y - centerY;
        const distance = Math.sqrt(dx * dx + dy * dy);
        const radius = this.wheelSize / 2 - 10;

        if (distance <= radius) {
            const angle = Math.atan2(dy, dx) * 180 / Math.PI;
            const hue = (angle + 360) % 360;
            const saturation = Math.min((distance / radius) * 100, 100);

            this.currentColor.hsv.h = hue;
            this.currentColor.hsv.s = saturation;

            this.updateFromHSV();
            this.updateIndicatorPosition();
        }

    }

    updateFromHSV() {
        if (this.isUpdatingFromUI) return;

        const {h, s, v} = this.currentColor.hsv;
        this.currentColor.rgb = this.hsvToRgb(h, s, v);
        this.currentColor.hsl = this.rgbToHsl(this.currentColor.rgb.r, this.currentColor.rgb.g, this.currentColor.rgb.b);
        this.currentColor.hex = this.rgbToHex(this.currentColor.rgb.r, this.currentColor.rgb.g, this.currentColor.rgb.b);

        this.updateUI();
        this.addToRecentColors(this.currentColor.hex);
    }

    updateFromRGB() {
        if (this.isUpdatingFromUI) return;
        this.isUpdatingFromUI = true;

        const r = Math.max(0, Math.min(255, parseInt(this.input.rgb.r.value) || 0));
        const g = Math.max(0, Math.min(255, parseInt(this.inputs.rgb.g.value) || 0));
        const b = Math.max(0, Math.min(255, parseInt(this.inputs.rgb.b.value) || 0));

        this.currentColor.rgb = {r, g, b};
        this.currentColor.hsv = this.rgbToHsv(r, g, b);
        this.currentColor.hsl = this.rgbToHsl(r, g, b);
        this.currentColor.hex = this.rgbToHex(r, g, b);

        this.updateUI();
        this.updateIndicatorPosition();
        this.addToRecentColors(this.currentColor.hex);

        this.isUpdatingFromUI = false;
    }

    updateFromHSL() {
        if (this.isUpdatingFromUI) return;
        this.isUpdatingFromUI = true;

        const h = Math.max(0, Math.min(360, parseInt(this.inputs.hsl.h.value) || 0));
        const s = Math.max(0, Math.min(100, parseInt(this.inputs.hsl.s.value) || 0));
        const l = Math.max(0, Math.min(100, parseInt(this.inputs.hsl.l.value) || 0));

        this.currentColor.hsl = {h, s, l};
        this.currentColor.rgb = this.hslToRgb(h, s, l);
        this.currentColor.hsv = this.rgbToHsv(this.currentColor.rgb.r, this.currentColor.rgb.g, this.currentColor.rgb.b);

        this.updateUI();
        this.updateIndicatorPosition();
        this.addToRecentColors(this.currentColor.hex);

        this.isUpdatingFromUI = false;
    }

    updateFromHex() {
        if (this.isUpdatingFromUI) return;
        this.isUpdatingFromUI = true;

        const hex = this.inputs.hex.value;
        const rgb = this.hexToRgb(hex);

        if (rgb) {
            this.currentColor.hex = hex;
            this.currentColor.rgb = rgb;
            this.currentColor.hsv = this.rgbToHsv(rgb.r, rgb.g, rgb.b);
            this.currentColor.hsl = this.rgbToHsl(rgb.r, rgb.g, rgb.b);

            this.updateUI();
            this.updateIndicatorPosition();
            this.addToRecentColors(hex);
        }

        this.isUpdatingFromUI = false;
    }

    updateUI() {
        if (this.isUpdatingFromUI) return;

        //Update color preview
        this.colorPreview.style.backgroundColor = this.currentColor.hex;

        //Update brightness slider
        this.brightnessSlider.value = this.currentColor.hsv.v;

        //Update the input fields (no event triggering)
        this.inputs.hex.value = this.currentColor.hex;

        this.inputs.rgb.r.value = this.currentColor.rgb.r;
        this.inputs.rgb.g.value = this.currentColor.rgb.g;
        this.inputs.rgb.b.value = this.currentColor.rgb.b;

        this.inputs.hsl.h.value = this.currentColor.hsl.h;
        this.inputs.hsl.s.value = this.currentColor.hsl.s;
        this.inputs.hsl.l.value = this.currentColor.hsl.l;

        this.inputs.hsv.h.value = Math.round(this.currentColor.hsv.h);
        this.inputs.hsv.s.value = Math.round(this.currentColor.hsv.s);
        this.inputs.hsv.v.value = Math.round(this.currentColor.hsv.v);
    }

    updateIndicatorPosition() {
        const centerX = this.wheelSize / 2;
        const centerY = this.wheelSize / 2;
        const radius = (this.wheel / 2 - 10) * (this.currentColor.hsv.s / 100);
        const angle = (this.currentColor.hsv.h * Math.PI) / 180;

        const x = centerX + radius * Math.cos(angle);
        const y = centerY + radius * Math.sin(angle);

        this.indicator.style.left = x + 'px';
        this.indicator.style.top = y + 'px';
    }

    //Storage management
    addToRecentColors(hex) {
        let recentColors = JSON.parse(localStorage.getItem('recentColors') || '[]');

        //Remove if already exists
        recentColors = recentColors.filter(color => color !== hex);

        //Add to beginning
        recentColors.unshift(hex);

        //Limit to maxRecentColors
        recentColors = recentColors.slice(0, this.maxRecentColors);

        localStorage.set.Item('recentColors', JSON.stringify(recentColors));
        this.renderRecentColors();
    }

    saveCurrentColor() {
        let savedColors = JSON.parse(localStorage.getItem('savedColors') || '[]');

        if (!savedColors.includes(this.currentColor.hex) && savedColors.length < this.maxSavedColors) {
            savedColors.push(this.currentColor.hex);
            localStorage.setItem('savedColors',JSON.stringify(savedColors));
            this.renderSavedPalette();
            this.showNotification('Color saved to palette!');
        } else if (savedColors.includes(this.currentColor.hex)) {
            this.showNotification('Color already in palette!');
        }else {
            this.showNotification('Palette is full!');
        }
    }

    
}
