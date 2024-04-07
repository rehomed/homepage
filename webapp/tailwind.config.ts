import type { Config } from "tailwindcss";
import daisyui from "daisyui";
import daisyThemes from "daisyui/src/theming/themes"

export default {
  content: ["./src/**/*.{ts,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [daisyui],
  daisyui: {
    themes: [
      {
        def: {
          ...daisyThemes.business,
          "--btn-focus-scale": "1",
        }
      }
    ],
  },
} satisfies Config;
