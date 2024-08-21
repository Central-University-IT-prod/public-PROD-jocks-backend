/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{html,js,,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "#5538ee",
        gray: { 100: "#DBDBDB", 200: "#6F6F6F" },
      },
    },
  },
  plugins: [],
};
