from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any, Iterable

from fastapi import APIRouter, Request


@dataclass(frozen=True)
class RouteMeta:
    title: str
    requires_auth: bool = False
    guest_only: bool = False
    roles: tuple[str, ...] = ()
    permissions: tuple[str, ...] = ()
    tags: tuple[str, ...] = ()

    def as_dict(self) -> dict[str, Any]:
        return {
            "title": self.title,
            "requires_auth": self.requires_auth,
            "guest_only": self.guest_only,
            "roles": list(self.roles),
            "permissions": list(self.permissions),
            "tags": list(self.tags),
        }


@dataclass(frozen=True)
class RouteDefinition:
    path: str
    name: str
    meta: RouteMeta
    methods: tuple[str, ...] = ("GET",)
    source: str = "frontend"
    summary: str | None = None
    description: str | None = None

    @property
    def normalized_path(self) -> str:
        normalized = self.path
        normalized = normalized.replace(":", "{")
        normalized = normalized.replace("/*", "/{path:path}")
        normalized = normalized.replace("{", "{")
        while "{" in normalized and "}" not in normalized:
            pass
        if "{" in normalized and "}" in normalized:
            normalized = normalized.replace("/", "/")
        path_parts = []
        for part in normalized.split("/"):
            if part.startswith("{") and part.endswith("}"):
                path_parts.append(part)
            elif part.endswith("}"):
                path_parts.append(part)
            elif part.startswith("{"):
                path_parts.append(part)
            else:
                path_parts.append(part)
        final_path = "/" + "/".join(part for part in path_parts if part or part == "")
        if final_path == "":
            return "/"

        final_path = final_path.replace("{id}", "{id}")
        final_path = final_path.replace("{path:path}", "{path:path}")
        final_path = final_path.replace("{/", "/")
        while "{{" in final_path:
            final_path = final_path.replace("{{", "{")
        while "}}" in final_path:
            final_path = final_path.replace("}}", "}")
        final_path = final_path.replace("//", "/")
        return final_path

    def body(self, request: Request) -> dict[str, Any]:
        route_path = request.url.path
        return {
            "name": self.name,
            "path": self.path,
            "route": route_path,
            "title": self.meta.title,
            "requires_auth": self.meta.requires_auth,
            "guest_only": self.meta.guest_only,
            "roles": list(self.meta.roles),
            "permissions": list(self.meta.permissions),
            "tags": list(self.meta.tags),
            "path_params": dict(request.path_params),
            "query_params": dict(request.query_params),
            "method": request.method,
            "source": self.source,
        }


def normalize_route_path(path: str) -> str:
    normalized = path.strip()
    if not normalized:
        return "/"
    normalized = normalized.replace("*path", "{path:path}")
    normalized = normalized.replace("*", "{wildcard:path}")
    normalized = normalized.replace(":id", "{id}")
    normalized = normalized.replace(":type", "{type}")
    normalized = normalized.replace(":busid", "{busid}")
    normalized = normalized.replace(":depotid", "{depotid}")
    normalized = normalized.replace(":tag", "{tag}")
    normalized = normalized.replace(":guildId", "{guildId}")
    normalized = normalized.replace(":channelId", "{channelId}")
    normalized = normalized.replace(":name", "{name}")
    normalized = normalized.replace(":slug", "{slug}")
    normalized = normalized.replace("{wildcard:path}", "{path:path}")
    if not normalized.startswith("/"):
        normalized = "/" + normalized
    return normalized


core_routes: list[RouteDefinition] = [
    RouteDefinition(path="/", name="Home", meta=RouteMeta(title="Home"), source="core"),
    RouteDefinition(path="/home", name="Home", meta=RouteMeta(title="Home"), source="core"),
    RouteDefinition(path="/login", name="Auth", meta=RouteMeta(title="Login", guest_only=True), source="core"),
    RouteDefinition(path="/profile", name="MyProfile", meta=RouteMeta(title="My Profile", requires_auth=True), source="core"),
    RouteDefinition(path="/user/:id", name="UserProfile", meta=RouteMeta(title="User Profile"), source="core"),
    RouteDefinition(path="/settings", name="Settings", meta=RouteMeta(title="Account Settings", requires_auth=True), source="core"),
    RouteDefinition(path="/map", name="MapGTA", meta=RouteMeta(title="Map"), source="core"),
    RouteDefinition(path="/cart", name="Cart", meta=RouteMeta(title="Shopping Cart", requires_auth=True), source="core"),
    RouteDefinition(path="/my-orders", name="MyOrders", meta=RouteMeta(title="My Orders", requires_auth=True), source="core"),
    RouteDefinition(path="/deliveries", name="Deliveries", meta=RouteMeta(title="Deliveries", requires_auth=True), source="core"),
    RouteDefinition(path="/delivery/create", name="Createdelivery", meta=RouteMeta(title="Create Delivery", requires_auth=True), source="core"),
    RouteDefinition(path="/delivery/track/:id", name="TrackDelivery", meta=RouteMeta(title="Track Delivery", requires_auth=True), source="core"),
    RouteDefinition(path="/delivery/:id", name="Delivery", meta=RouteMeta(title="Delivery", requires_auth=True), source="core"),
    RouteDefinition(path="/dash/driver", name="DriverDash", meta=RouteMeta(title="Driver Dashboard", requires_auth=True, roles=("driver", "admin")), source="core"),
    RouteDefinition(path="/wallet", name="Wallet", meta=RouteMeta(title="Wallet", requires_auth=True), source="core"),
]

legal_routes: list[RouteDefinition] = [
    RouteDefinition(path="/about", name="About", meta=RouteMeta(title="About Us"), source="legal"),
    RouteDefinition(path="/contact", name="Contact", meta=RouteMeta(title="Contact Us"), source="legal"),
    RouteDefinition(path="/faq", name="Faq", meta=RouteMeta(title="FAQ"), source="legal"),
    RouteDefinition(path="/terms", name="Terms", meta=RouteMeta(title="Terms of Service"), source="legal"),
    RouteDefinition(path="/privacy", name="Privacy", meta=RouteMeta(title="Privacy Policy"), source="legal"),
    RouteDefinition(path="/refund", name="Refund", meta=RouteMeta(title="Refund Policy"), source="legal"),
    RouteDefinition(path="/shipping", name="Shipping", meta=RouteMeta(title="Shipping Information"), source="legal"),
    RouteDefinition(path="/returns", name="Returns", meta=RouteMeta(title="Returns Policy"), source="legal"),
    RouteDefinition(path="/disclaimer", name="Disclaimer", meta=RouteMeta(title="Disclaimer"), source="legal"),
    RouteDefinition(path="/blog", name="Blog", meta=RouteMeta(title="Blog"), source="legal"),
]

error_routes: list[RouteDefinition] = [
    RouteDefinition(path="/404", name="NotFound", meta=RouteMeta(title="Page Not Found"), source="error"),
    RouteDefinition(path="/error/404", name="NotFound", meta=RouteMeta(title="Page Not Found"), source="error"),
    RouteDefinition(path="/403", name="Forbidden", meta=RouteMeta(title="Access Denied"), source="error"),
    RouteDefinition(path="/error/403", name="Forbidden", meta=RouteMeta(title="Access Denied"), source="error"),
]

admin_routes: list[RouteDefinition] = [
    RouteDefinition(path="/admin", name="Admin", meta=RouteMeta(title="Admin", requires_auth=True, roles=("admin",)), source="admin"),
    RouteDefinition(path="/admin/role-requests", name="RoleRequestsPage", meta=RouteMeta(title="Role Requests", requires_auth=True, roles=("admin",)), source="admin"),
    RouteDefinition(path="/admin/moderator-applications", name="ModeratorApplicationsPage", meta=RouteMeta(title="Moderator Applications", requires_auth=True, roles=("admin",)), source="admin"),
    RouteDefinition(path="/dash", name="Dash", meta=RouteMeta(title="Dashboard", requires_auth=True), source="admin"),
]

farms_routes: list[RouteDefinition] = [
    RouteDefinition(path="/dash", name="Dash", meta=RouteMeta(title="Dashboard", requires_auth=True, roles=("farmer", "admin")), source="farms"),
    RouteDefinition(path="/farms", name="Farms", meta=RouteMeta(title="Farms"), source="farms"),
    RouteDefinition(path="/create-farm", name="CreateFarm", meta=RouteMeta(title="Create Farm", requires_auth=True), source="farms"),
    RouteDefinition(path="/tools", name="Tools", meta=RouteMeta(title="Tools"), source="farms"),
    RouteDefinition(path="/products", name="Products", meta=RouteMeta(title="Products"), source="farms"),
    RouteDefinition(path="/recipes", name="Recipes", meta=RouteMeta(title="Recipes"), source="farms"),
    RouteDefinition(path="/create-recipe", name="CreateRecipe", meta=RouteMeta(title="Create Recipe", requires_auth=True), source="farms"),
    RouteDefinition(path="/products/:type/:id", name="Product", meta=RouteMeta(title="Product"), source="farms"),
    RouteDefinition(path="/farm/:id", name="Farm", meta=RouteMeta(title="Farm"), source="farms"),
    RouteDefinition(path="/recipe/:id", name="Recipe", meta=RouteMeta(title="Recipe"), source="farms"),
]

bus_routes: list[RouteDefinition] = [
    RouteDefinition(path="/bus", name="BusHome", meta=RouteMeta(title="Bus Home"), source="bus"),
    RouteDefinition(path="/bus/time", name="BusTime", meta=RouteMeta(title="Bus Timetable"), source="bus"),
    RouteDefinition(path="/bus/map", name="BusMap", meta=RouteMeta(title="Bus Map"), source="bus"),
    RouteDefinition(path="/bus/time/:busid", name="BusRouteDetail", meta=RouteMeta(title="Bus Route Details"), source="bus"),
    RouteDefinition(path="/bus/live/:busid", name="BusLiveTracking", meta=RouteMeta(title="Live Bus Tracking"), source="bus"),
    RouteDefinition(path="/bus/depot/:depotid", name="BusDepot", meta=RouteMeta(title="Bus Depot"), source="bus"),
]

crop_routes: list[RouteDefinition] = [
    RouteDefinition(path="/crops", name="Crops", meta=RouteMeta(title="Crops"), source="crops"),
    RouteDefinition(path="/bazarbhav", name="BazaarBhav", meta=RouteMeta(title="Bazaar Bhav"), source="crops"),
    RouteDefinition(path="/grocery", name="Crops", meta=RouteMeta(title="Grocery"), source="crops"),
    RouteDefinition(path="/crop/:id", name="Crop", meta=RouteMeta(title="Crop"), source="crops"),
    RouteDefinition(path="/aboutcrop/:id", name="AboutCrop", meta=RouteMeta(title="About Crop"), source="crops"),
]

events_routes: list[RouteDefinition] = [
    RouteDefinition(path="/events", name="Events", meta=RouteMeta(title="Events"), source="events"),
    RouteDefinition(path="/create-event", name="CreateEvent", meta=RouteMeta(title="Create Event", requires_auth=True), source="events"),
    RouteDefinition(path="/artists", name="Artists", meta=RouteMeta(title="Artists"), source="events"),
    RouteDefinition(path="/create-artist", name="CreateArtist", meta=RouteMeta(title="Create Artist", requires_auth=True), source="events"),
    RouteDefinition(path="/vendors", name="Vendors", meta=RouteMeta(title="Vendors"), source="events"),
    RouteDefinition(path="/music", name="Music", meta=RouteMeta(title="Music"), source="events"),
    RouteDefinition(path="/event/:id/tickets", name="EventTickets", meta=RouteMeta(title="Event Tickets", requires_auth=True), source="events"),
    RouteDefinition(path="/event/:id", name="Event", meta=RouteMeta(title="Event"), source="events"),
    RouteDefinition(path="/artist/:id", name="Artist", meta=RouteMeta(title="Artist"), source="events"),
]

baito_routes: list[RouteDefinition] = [
    RouteDefinition(path="/baitos", name="Baitos", meta=RouteMeta(title="Baitos"), source="baito"),
    RouteDefinition(path="/baitos/dash", name="BaitoDash", meta=RouteMeta(title="Baito Dashboard", requires_auth=True, roles=("worker", "admin")), source="baito"),
    RouteDefinition(path="/baitos/hire", name="HireWorkers", meta=RouteMeta(title="Hire Workers"), source="baito"),
    RouteDefinition(path="/baitos/create-profile", name="CreateBaitoProfile", meta=RouteMeta(title="Create Baito Profile", requires_auth=True, roles=("worker", "admin")), source="baito"),
    RouteDefinition(path="/create-baito", name="CreateBaito", meta=RouteMeta(title="Create Baito", requires_auth=True), source="baito"),
    RouteDefinition(path="/booking", name="Booking", meta=RouteMeta(title="Booking"), source="baito"),
    RouteDefinition(path="/baito/:id", name="Baito", meta=RouteMeta(title="Baito"), source="baito"),
    RouteDefinition(path="/baitos/worker/:id", name="Worker", meta=RouteMeta(title="Worker", requires_auth=True), source="baito"),
]

social_routes: list[RouteDefinition] = [
    RouteDefinition(path="/posts", name="Posts", meta=RouteMeta(title="Posts"), source="social"),
    RouteDefinition(path="/create-post", name="CreatePost", meta=RouteMeta(title="Create Post", requires_auth=True), source="social"),
    RouteDefinition(path="/post/:id", name="Post", meta=RouteMeta(title="Post"), source="social"),
    RouteDefinition(path="/hashtag/:tag", name="Hashtag", meta=RouteMeta(title="Hashtag"), source="social"),
]

chats_routes: list[RouteDefinition] = [
    RouteDefinition(path="/merechats", name="MeChats", meta=RouteMeta(title="My Chats", requires_auth=True), source="chats"),
    RouteDefinition(path="/newchats", name="NewChats", meta=RouteMeta(title="New Chats", requires_auth=True), source="chats"),
    RouteDefinition(path="/merechats/:id", name="OneChatPage", meta=RouteMeta(title="Chat", requires_auth=True), source="chats"),
    RouteDefinition(path="/newchat/:id", name="NewChatPage", meta=RouteMeta(title="Chat", requires_auth=True), source="chats"),
]

places_routes: list[RouteDefinition] = [
    RouteDefinition(path="/map", name="MapGTA", meta=RouteMeta(title="Map"), source="places"),
    RouteDefinition(path="/places", name="Places", meta=RouteMeta(title="Places"), source="places"),
    RouteDefinition(path="/create-place", name="CreatePlace", meta=RouteMeta(title="Create Place", requires_auth=True), source="places"),
    RouteDefinition(path="/itinerary", name="Itinerary", meta=RouteMeta(title="Itinerary"), source="places"),
    RouteDefinition(path="/create-itinerary", name="CreateItinerary", meta=RouteMeta(title="Create Itinerary", requires_auth=True), source="places"),
    RouteDefinition(path="/edit-itinerary", name="EditItinerary", meta=RouteMeta(title="Edit Itinerary", requires_auth=True), source="places"),
    RouteDefinition(path="/place/:id", name="Place", meta=RouteMeta(title="Place"), source="places"),
    RouteDefinition(path="/itinerary/:id", name="Itinerary", meta=RouteMeta(title="Itinerary"), source="places"),
]

route_groups: dict[str, Iterable[RouteDefinition]] = {
    "core": core_routes,
    "legal": legal_routes,
    "error": error_routes,
    "admin": admin_routes,
    "farms": farms_routes,
    "bus": bus_routes,
    "crops": crop_routes,
    "events": events_routes,
    "baito": baito_routes,
    "social": social_routes,
    "chats": chats_routes,
    "places": places_routes,
}


def merge_routes() -> list[RouteDefinition]:
    merged: list[RouteDefinition] = []
    seen_paths: set[str] = set()
    for group in route_groups.values():
        for route in group:
            if route.path in seen_paths:
                continue
            seen_paths.add(route.path)
            merged.append(route)
    return merged


route_registry = merge_routes()
router = APIRouter(prefix="", tags=["routes"])


def build_route_handler(route: RouteDefinition):
    async def handler(request: Request):
        return route.body(request)

    handler.__name__ = route.name
    return handler


for route in route_registry:
    normalized = normalize_route_path(route.path)
    router.add_api_route(
        path=normalized,
        endpoint=build_route_handler(route),
        methods=list(route.methods),
        summary=route.summary or route.meta.title,
        description=route.description or f"Frontend route for {route.name}.",
        tags=list(route.meta.tags) or [route.source],
    )


__all__ = ["RouteDefinition", "RouteMeta", "route_registry", "router", "normalize_route_path"]
