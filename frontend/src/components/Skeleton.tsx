interface SkeletonProps {
  width?: string | number;
  height?: string | number;
  radius?: string | number;
  className?: string;
  style?: React.CSSProperties;
}

export function Skeleton({ width, height, radius, className, style }: SkeletonProps) {
  return (
    <div
      className={`skeleton ${className ?? ""}`}
      style={{
        width: width ?? "100%",
        height: height ?? "1em",
        borderRadius: radius ?? 8,
        ...style,
      }}
    />
  );
}

export function SkeletonCard() {
  return (
    <div className="skeleton-card">
      <Skeleton height={140} radius={12} />
      <Skeleton width="60%" height={14} style={{ marginTop: 12 }} />
      <Skeleton width="90%" height={12} style={{ marginTop: 8 }} />
      <Skeleton width="40%" height={12} style={{ marginTop: 8 }} />
    </div>
  );
}

export function SkeletonGrid({ count = 6 }: { count?: number }) {
  return (
    <div className="cards-grid">
      {Array.from({ length: count }).map((_, i) => (
        <SkeletonCard key={i} />
      ))}
    </div>
  );
}

export function SkeletonList({ count = 4 }: { count?: number }) {
  return (
    <div style={{ display: "grid", gap: "0.75rem" }}>
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="skeleton-list-row">
          <Skeleton width={80} height={80} radius={12} />
          <div style={{ flex: 1, display: "grid", gap: 8 }}>
            <Skeleton width="40%" height={14} />
            <Skeleton width="70%" height={12} />
            <Skeleton width="55%" height={12} />
          </div>
        </div>
      ))}
    </div>
  );
}
