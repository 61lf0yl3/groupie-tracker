// Initialize Communication with Back-end Services
var platform = new H.service.Platform({
    apikey: "Z30Ae5EJvVXg_yeLJo67KsyoTp2rP0v3dQCQHV0bfFY"
});

// Obtain the default map types from the platform object:
var defaultLayers = platform.createDefaultLayers();

// Initialize a map - this map is centered over Europe
var map = new H.Map(document.getElementById('map'),
    defaultLayers.vector.normal.map, {
    center: { lat: 48.499998, lng: 23.3833318 },
    zoom: 2.5,
    pixelRatio: window.devicePixelRatio || 1
});

// Add a resize listener to make sure that the map occupies the whole container
window.addEventListener('resize', () => map.getViewPort().resize());

// Enable the event system on the map instance:
var mapEvents = new H.mapevents.MapEvents(map);

// Instantiate the default behavior, providing the mapEvents object:
new H.mapevents.Behavior(mapEvents);

// Create the default UI components
var ui = H.ui.UI.createDefault(map, defaultLayers);

// Read locations and add markers
  var cities = document.getElementsByClassName('cities');
  for (let i = 0; i < cities.length; i++) {
    // Create the parameters for the geocoding request:
    var geocodingParams = {
        searchText: cities[i].textContent
    };

    // Define a callback function to process the geocoding response:
    var onResult = function(result) {
    var locations = result.Response.View[0].Result,
        position,
        marker;

    // Add a marker for each location found
    position = {
        lat: locations[0].Location.DisplayPosition.Latitude,
        lng: locations[0].Location.DisplayPosition.Longitude
    };
    marker = new H.map.Marker(position);
    map.addObject(marker);
    };
    
    // Get an instance of the geocoding service:
    var geocoder = platform.getGeocodingService();
    
    // Call the geocode method with the geocoding parameters,
    // the callback and an error callback function (called if a
    // communication error occurs):
    geocoder.geocode(geocodingParams, onResult, function(e) {
    alert(e);
    });
  }



    