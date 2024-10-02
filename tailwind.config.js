/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/html/**/*.{html,js,tmpl}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Inter var", ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
