import type { Config } from 'tailwindcss'

export default {
    content: [
        './app/**/*.{vue,ts,js}',
        './components/**/*.{vue,ts,js}',
        './layouts/**/*.{vue,ts,js}',
        './pages/**/*.{vue,ts,js}',
    ],
    theme: {
        extend: {
            colors: {
                navy: '#1E3A8A',
                gold: '#CA8A04',
                action: '#D31245',
            },
            fontFamily: {
                heading: ['"Playfair Display SC"', 'serif'],
                body: ['Karla', '"Microsoft YaHei"', 'sans-serif'],
            },
        },
    },
    plugins: [],
} satisfies Config
