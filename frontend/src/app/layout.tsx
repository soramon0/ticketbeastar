import { Inter } from 'next/font/google';
import './globals.css';
import { classNames } from '@/utils';

export const metadata = {
  title: 'Ticket Beastar',
  description: 'Buy a concert ticket easily',
};

const inter = Inter({ subsets: ['latin'] });

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={classNames(inter.className, 'bg-gray-50')}>
        {children}
      </body>
    </html>
  );
}
