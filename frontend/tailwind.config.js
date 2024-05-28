/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{html,ts}"],

  theme: {
    extend: {
      boxShadow: {
        menu: "-4px 0 20px 0px rgba(0, 0, 0, 0.85)",
        "light-shadow": "0 0 12px 0px rgba(120, 120, 120, 0.06)",
      },
    },
  },
  plugins: [],
};
