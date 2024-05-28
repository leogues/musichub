export const formatNumber = (number: number): string => {
  const units = ["", "K", "m"];
  const div = 1000;

  for (let i = 0; i < units.length; i++) {
    if (number < div)
      return (
        number.toLocaleString("pt-BR", { maximumFractionDigits: 2 }) + units[i]
      );
    number = number / div;
  }

  return number.toString();
};
