export interface MapDimensions {
    width: number;
    height: number;
}

export interface MapAsset {
    image: string;
    fallbackImage?: string;
    dimensions?: MapDimensions;
}

export interface Point2D {
    x: number;
    y: number;
}

export interface LiveEvent {
    isLive?: boolean;
    eventName?: string;
    endsAt?: string | number | Date | null;
    remainingSecs?: number;
    rewardMultiplier?: number;
}

export interface LocationDetails {
    address?: string;
    price?: number;
    intelData?: string;
}

export interface GtaFloor {
    level: number;
    name?: string;
    image?: string;
    locations?: GtaLocation[];
}

export interface GtaLocation {
    id?: string;
    name: string;
    description?: string;
    category: string;
    x: number;
    y: number;
    floorLevel?: number;
    icon?: string;
    iconUrl?: string;
    details?: LocationDetails;
    liveEvent?: LiveEvent;
    membersOnly?: boolean;
    // Compatibility aliases used by the UI layer
    locationId?: string;
    [key: string]: unknown;
}

export interface LiveEntity {
    id: string;
    type: "vehicle" | "player" | string;
    name: string;
    position: Point2D;
    heading?: number;
    speed?: number;
    floor?: number;
    occupants?: number;
    updated?: string | number | Date;
}

export interface TerritoryPoint {
    x: number;
    y: number;
}

export interface Territory {
    id?: string;
    name: string;
    gangName?: string;
    owner?: string;
    controlPct: number;
    color?: string;
    polygonPoints?: TerritoryPoint[];
    points?: TerritoryPoint[];
    [key: string]: unknown;
}

export interface CategoryFilter {
    id: string;
    label: string;
    icon?: string;
    count?: number;
}

export type CategoryItem = CategoryFilter;

export interface LockedArea {
    id?: string;
    label?: string;
    x: number;
    y: number;
    width?: number;
    height?: number;
    condition?: string;
    dependsOn?: string;
}

export interface PermalinkInfo {
    url?: string;
    targetId?: string;
    entity?: string;
    zoom?: number;
    focusPoint?: Point2D;
    floorLevel?: number | null;
}

export interface MapConfig {
    entity?: string;
    title?: string;
    map?: MapAsset;
    lockedAreas?: LockedArea[];
    floors?: GtaFloor[];
    [key: string]: unknown;
}

export interface MapResponseData {
    entity?: string;
    title?: string;
    map?: MapAsset;
    categories?: CategoryFilter[];
    locations?: GtaLocation[];
    territories?: Territory[];
    lockedAreas?: LockedArea[];
    floors?: GtaFloor[];
    playerProgress?: Record<string, number>;
    permalink?: PermalinkInfo;
    [key: string]: unknown;
}

export interface GtaMapState {
    zoomLevel: number;
    panX: number;
    panY: number;
    angle: number;
    flip: boolean;
    isDragging: boolean;
    startX: number;
    startY: number;
    velocityX: number;
    velocityY: number;
    currentIndex: number;
    activeEntity: string;

    floors: GtaFloor[];
    currentFloor: number | null;

    locations: GtaLocation[];
    activeCategories: Set<string>;
    liveEntities: Map<string, LiveEntity>;
    customWaypoint: TerritoryPoint | null;
    activeMission: { from: TerritoryPoint; to: TerritoryPoint } | null;
    deliveryMissions: any[];
    territories: Territory[];
    liveEvents: any[];

    isMeasuring: boolean;
    measurePoints: TerritoryPoint[];
    cursorCoords: TerritoryPoint;
    timerIntervals: number[];
    wsConnection: WebSocket | null;
    reconnectTimer: number | null;
    baseMapImageSrc?: string;
}
