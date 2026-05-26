const petals = Array.from({ length: 24 }, (_, index) => ({
  id: index,
  left: 8 + ((index * 37) % 90),
  size: 4 + (index % 7),
  opacity: 0.25 + (index % 5) * 0.08,
  delay: -(index * 1.35),
  duration: 13 + (index % 9) * 1.45,
  drift: 80 + (index % 6) * 24,
  rotate: (index * 23) % 180,
}));

export function SakuraPetals() {
  return (
    <div className="sakura-petal-layer pointer-events-none absolute inset-0 z-[3] overflow-hidden" aria-hidden="true">
      {petals.map((petal) => (
        <span
          key={petal.id}
          className="sakura-petal absolute block rounded-[65%_35%_70%_30%] bg-rose-100/80 shadow-[0_0_10px_rgba(255,255,255,0.22)]"
          style={{
            left: `${petal.left}%`,
            width: `${petal.size}px`,
            height: `${Math.max(4, petal.size - 1)}px`,
            opacity: petal.opacity,
            animationDelay: `${petal.delay}s`,
            animationDuration: `${petal.duration}s`,
            ["--petal-drift" as string]: `${petal.drift}px`,
            ["--petal-rotate" as string]: `${petal.rotate}deg`,
          }}
        />
      ))}
    </div>
  );
}
