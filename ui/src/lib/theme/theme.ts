import { Button, createTheme } from "@mantine/core";
import classes from './theme.module.css';

export const someTheme = createTheme({
  fontFamily: 'Funnel Display, sans-serif',
  primaryColor: "red",
  colors: {
    someRed: [
      "#ffeae9",
      "#ffd6d1",
      "#f7aba4",
      "#f17e73",
      "#eb584a",
      "#e83f2f",
      "#E33226",
      "#cd2315",
      "#BB281D",
      "#a1100a"
    ],
    someYellow: [
      "#fffbeb",
      "#fff3cb",
      "#ffe488",
      "#ffd24f",
      "#ffbc20",
      "#f99a07",
      "#dd7202",
      "#b74f06",
      "#943c0c",
      "#461902",
    ],
    someBlue: [
      "#f0f9ff",
      "#e1f2fd",
      "#bce5fb",
      "#81d1f8",
      "#3dbaf3",
      "#1398d6",
      "#0881c1",
      "#07679d",
      "#0b5781",
      "#0f486b",
      "#0a2f47",
    ],
  },
  // Same breakpoint values as tailwind
  breakpoints: {
    xs: "36em",
    sm: "40em",
    md: "48em",
    lg: "64em",
    xl: "80em",
    xxl: "96em",
    xxxl: "142em",
    xxxxl: "172em",
  },
  components: {
    Button: Button.extend({ classNames: classes }),
  }
});
