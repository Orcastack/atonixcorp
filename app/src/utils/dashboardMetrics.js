import { useEffect, useState } from 'react';

const clamp = (value, min, max) => Math.min(max, Math.max(min, value));

export const useAnimatedNumber = (target, duration = 900) => {
  const [value, setValue] = useState(() => Number(target) || 0);

  useEffect(() => {
    const startValue = Number(value) || 0;
    const endValue = Number(target) || 0;
    const delta = endValue - startValue;
    const startTime = window.performance.now();
    let frameId = 0;

    const step = (timestamp) => {
      const progress = clamp((timestamp - startTime) / duration, 0, 1);
      const eased = 1 - Math.pow(1 - progress, 3);
      setValue(startValue + delta * eased);
      if (progress < 1) {
        frameId = window.requestAnimationFrame(step);
      }
    };

    frameId = window.requestAnimationFrame(step);
    return () => window.cancelAnimationFrame(frameId);
  }, [duration, target, value]);

  return value;
};

export const formatAnimatedMetric = (value, suffix = '') => {
  if (Number.isInteger(value)) {
    return `${value}${suffix}`;
  }
  return `${value.toFixed(1)}${suffix}`;
};

export const buildBalancedMetricOrder = (items, seed = 0) => {
  const list = [...items];
  if (list.length < 3) return list;

  const pivot = Math.abs(Number(seed) || 0) % list.length;
  return [...list.slice(pivot), ...list.slice(0, pivot)].map((item, index) => ({
    ...item,
    metricClassName: `${item.metricClassName || ''} premium-metric-card--slot-${(index % 4) + 1}`.trim(),
  }));
};
