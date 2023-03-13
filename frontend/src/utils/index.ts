export function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ');
}

interface Options {
  amount: number | string;
  locale?: string;
  currency?: string;
}

export function formatCurrency({ amount, locale = 'en', currency }: Options) {
  const options: Intl.NumberFormatOptions = {
    minimumFractionDigits: 2,
  };

  if (currency) {
    options.currency = currency;
    options.style = 'currency';
  }

  const fn = new Intl.NumberFormat(locale, options);
  return fn.format(Number(amount));
}
