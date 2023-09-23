/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./pkg/web/views/**/*.{html,js}"],
  theme: {
    extend: {
      colors: {
        primary: "#0f172a",
      },
      fontFamily: {
        Kanit: ["Kanit, sans-serif"],
      },
    },
  },
  plugins: [],
}

