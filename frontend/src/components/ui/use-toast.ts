import { useState, useCallback } from 'react';

interface ToastOptions {
  title: string;
  description?: string;
  variant?: 'default' | 'destructive';
  duration?: number;
}

export function useToast() {
  const [toasts, setToasts] = useState<ToastOptions[]>([]);

  const toast = useCallback(({ title, description, variant = 'default', duration = 3000 }: ToastOptions) => {
    const id = Date.now();
    setToasts(prev => [...prev, { title, description, variant }]);
    setTimeout(() => {
      setToasts(prev => prev.filter(toast => toast !== id));
    }, duration);
  }, []);

  return { toast, toasts };
} 