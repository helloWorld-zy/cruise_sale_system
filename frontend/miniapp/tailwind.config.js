/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
    "./pages/**/*.{vue,js,ts,jsx,tsx}",
    "./components/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#0891B2',
        secondary: '#22D3EE',
        cta: '#22C55E',
        background: '#ECFEFF',
        text: '#164E63'
      },
      fontFamily: {
        heading: ['Varela Round', 'sans-serif'],
        body: ['Nunito Sans', 'sans-serif']
      },
      boxShadow: {
        'sm': '0 1px 2px rgba(0,0,0,0.05)',
        'md': '0 4px 6px rgba(0,0,0,0.1)',
        'lg': '0 10px 15px rgba(0,0,0,0.1)',
        'xl': '0 20px 25px rgba(0,0,0,0.15)'
      }
    },
  },
  plugins: [],
}
