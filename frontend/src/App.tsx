import type { KonvaEventObject } from "konva/lib/Node";
import { useRef, useState } from "react";
import { Stage, Layer, Circle } from "react-konva";

interface Hold {
    x: number;
    y: number;
    id: number;
}

export default function App() {
    const [imageSrc, setImageSrc] = useState<string | null>(null);
    const [holds, setHolds] = useState<Hold[]>([]);
    const [saving, setSaving] = useState(false)
    const [savedId, setSavedId] = useState<string | null>(null)
    const [stageSize, setStageSize] = useState({ width: 0, height: 0 });
    const imageRef = useRef<HTMLImageElement>(null);
    const containerRef = useRef<HTMLDivElement>(null);

    const API = "http://localhost:8080"

    async function handleSave() {
        if (holds.length === 0) return
        setSaving(true)
        const id = await saveRoute(holds)
        setSavedId(id)
        setSaving(false)
    }

    async function saveRoute(holds: Hold[]) {
        // 1. Create the route
        const routeRes = await fetch(`${API}/api/routes`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name: "New Route", grade: "V?" }),
        })
        const route = await routeRes.json()

        // 2. Save the holds
        const holdsRes = await fetch(`${API}/api/routes/${route.id}/holds`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(
                holds.map((h, i) => ({
                    x_pct: h.x,
                    y_pct: h.y,
                    note: "",
                    seq_order: i + 1,
                }))
            ),
        })
        const saved = await holdsRes.json()
        console.log("saved route:", route.id, "holds:", saved)
        return route.id
    }

    function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0];
        if (!file) return;
        setHolds([]);
        setImageSrc(URL.createObjectURL(file));
    }

    function handleImageLoad() {
        const img = imageRef.current;
        if (!img) return;
        setStageSize({ width: img.clientWidth, height: img.clientHeight });
    }

    function handleStageClick(e: KonvaEventObject<MouseEvent> | KonvaEventObject<TouchEvent>) {
        const stage = e.target.getStage();
        if (!stage) return;
        const pos = stage.getPointerPosition();
        if (!pos) return;

        // Store as normalized 0.0–1.0 coords
        const hold: Hold = {
            id: Date.now(),
            x: pos.x / stageSize.width,
            y: pos.y / stageSize.height,
        };
        setHolds((prev) => [...prev, hold]);
    }

    return (
        <div style={{ maxWidth: 800, margin: "0 auto", padding: 24 }}>
            <h2>Boulder Beta</h2>

            <input
                type="file"
                accept="image/*"
                capture="environment"
                onChange={handleFileChange}
            />

            {imageSrc && (
                <div
                    ref={containerRef}
                    style={{ position: "relative", marginTop: 16 }}
                >
                    <img
                        ref={imageRef}
                        src={imageSrc}
                        onLoad={handleImageLoad}
                        style={{ width: "100%", display: "block" }}
                    />

                    {stageSize.width > 0 && (
                        <Stage
                            width={stageSize.width}
                            height={stageSize.height}
                            onClick={handleStageClick}
                            style={{ position: "absolute", top: 0, left: 0 }}
                        >
                            <Layer>
                                {holds.map((hold) => (
                                    <Circle
                                        key={hold.id}
                                        x={hold.x * stageSize.width}
                                        y={hold.y * stageSize.height}
                                        radius={14}
                                        stroke="#00ff88"
                                        strokeWidth={3}
                                        fill="rgba(0,255,136,0.15)"
                                    />
                                ))}
                            </Layer>
                        </Stage>
                    )}
                </div>
            )}

            {holds.length > 0 && (
                <p style={{ color: "#999", fontSize: 13, marginTop: 8 }}>
                    {holds.length} hold{holds.length > 1 ? "s" : ""} placed —{" "}
                    <span
                        style={{ cursor: "pointer", color: "#e94560" }}
                        onClick={() => setHolds([])}
                    >
                        clear all
                    </span>
                </p>
            )}

            {imageSrc && holds.length > 0 && (
                <button onClick={handleSave} disabled={saving}>
                    {saving ? "Saving..." : "Save Route"}
                </button>
            )}

            {savedId && (
                <p style={{ color: "#00ff88", fontSize: 13 }}>
                    Saved! Route ID: {savedId}
                </p>
            )}
        </div>
    );
}
