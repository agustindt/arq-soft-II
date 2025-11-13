// Seed data para MongoDB (Activities y Reservations)
// Este archivo contiene datos de prueba para development/testing
// Uso: docker-compose exec mongo mongosh < /seed.js

print("===========================================");
print("Seeding MongoDB - Activities & Reservations");
print("===========================================\n");

// ===== Activities Database =====
print("→ Seeding activitiesdb...");

db = db.getSiblingDB('activitiesdb');

// Limpiar colección de prueba (SOLO para testing!)
// db.activities.deleteMany({ name: /TEST/ });

// Actividades de prueba variadas
const activities = [
  {
    name: "Fútbol 5 Nocturno",
    description: "Partido de fútbol 5 en cancha sintética bajo techo. Ideal para después del trabajo. Incluye arbitraje y balón.",
    category: "football",
    difficulty: "intermediate",
    duration: 90,
    maxParticipants: 10,
    location: "Polideportivo Central, Cancha 3",
    price: 500.0,
    imageUrl: "https://example.com/futbol5.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Yoga Matutino en el Parque",
    description: "Sesión de yoga al aire libre para principiantes y avanzados. Trae tu mat. Instructor certificado.",
    category: "yoga",
    difficulty: "beginner",
    duration: 60,
    maxParticipants: 20,
    location: "Parque de la Reserva, Zona Verde",
    price: 250.0,
    imageUrl: "https://example.com/yoga.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Running Club 10K",
    description: "Entrenamiento grupal de running para preparación de 10K. Incluye plan de entrenamiento y análisis de técnica.",
    category: "running",
    difficulty: "advanced",
    duration: 120,
    maxParticipants: 30,
    location: "Malecón Costa Verde",
    price: 300.0,
    imageUrl: "https://example.com/running.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Básquetbol 3x3",
    description: "Torneo de básquetbol 3x3 estilo streetball. Premios para los ganadores.",
    category: "basketball",
    difficulty: "intermediate",
    duration: 120,
    maxParticipants: 12,
    location: "Cancha Municipal Norte",
    price: 400.0,
    imageUrl: "https://example.com/basketball.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Natación para Principiantes",
    description: "Clases de natación para adultos que están empezando. Incluye instructor y equipo básico.",
    category: "swimming",
    difficulty: "beginner",
    duration: 60,
    maxParticipants: 8,
    location: "Club Náutico, Piscina Olímpica",
    price: 450.0,
    imageUrl: "https://example.com/swimming.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Ciclismo de Montaña",
    description: "Ruta de mountain bike por cerros locales. Nivel intermedio-avanzado. Trae tu propia bici.",
    category: "cycling",
    difficulty: "advanced",
    duration: 180,
    maxParticipants: 15,
    location: "Inicio: Parque Zonal Huáscar",
    price: 350.0,
    imageUrl: "https://example.com/cycling.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Tenis Singles - Torneo Amateur",
    description: "Torneo de tenis singles amateur. Canchas de arcilla. Sistema de eliminación directa.",
    category: "tennis",
    difficulty: "intermediate",
    duration: 90,
    maxParticipants: 16,
    location: "Club de Tenis Los Alamos",
    price: 600.0,
    imageUrl: "https://example.com/tennis.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Volleyball Playa Mixto",
    description: "Partidos de volleyball en la playa. Equipos mixtos. Ambiente relajado y divertido.",
    category: "volleyball",
    difficulty: "beginner",
    duration: 120,
    maxParticipants: 12,
    location: "Playa Agua Dulce",
    price: 200.0,
    imageUrl: "https://example.com/volleyball.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "CrossFit Extremo",
    description: "WOD intensivo de CrossFit. Solo para atletas avanzados. Incluye coach certificado.",
    category: "fitness",
    difficulty: "advanced",
    duration: 90,
    maxParticipants: 15,
    location: "Box CrossFit Revolution",
    price: 550.0,
    imageUrl: "https://example.com/crossfit.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Paddle Tennis Dobles",
    description: "Sesión de pádel en modalidad dobles. Todos los niveles bienvenidos. Raquetas disponibles.",
    category: "paddle",
    difficulty: "beginner",
    duration: 90,
    maxParticipants: 8,
    location: "Centro Deportivo Las Palmas",
    price: 400.0,
    imageUrl: "https://example.com/paddle.jpg",
    createdAt: new Date(),
    updatedAt: new Date()
  }
];

const activityResult = db.activities.insertMany(activities);
print(`✓ Insertadas ${activityResult.insertedIds.length} actividades en activitiesdb`);

// Mostrar actividades por categoría
print("\n→ Actividades por categoría:");
db.activities.aggregate([
  { $group: { _id: "$category", count: { $sum: 1 } } },
  { $sort: { count: -1 } }
]).forEach(doc => {
  print(`  - ${doc._id}: ${doc.count} actividades`);
});

// ===== Reservations Database =====
print("\n→ Seeding reservasdb...");

db = db.getSiblingDB('reservasdb');

// Obtener IDs de actividades para las reservas
db = db.getSiblingDB('activitiesdb');
const activityIds = db.activities.find({}, { _id: 1 }).limit(5).toArray().map(a => a._id.toString());

db = db.getSiblingDB('reservasdb');

// Reservas de prueba
const reservations = [
  {
    activityId: activityIds[0] || "sample-activity-id-1",
    userId: "johndoe",
    userName: "John Doe",
    userEmail: "john.doe@test.com",
    date: new Date("2025-12-15T18:00:00Z"),
    participants: 5,
    status: "confirmed",
    totalPrice: 500.0,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    activityId: activityIds[1] || "sample-activity-id-2",
    userId: "janesmth",
    userName: "Jane Smith",
    userEmail: "jane.smith@test.com",
    date: new Date("2025-12-16T07:00:00Z"),
    participants: 1,
    status: "confirmed",
    totalPrice: 250.0,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    activityId: activityIds[2] || "sample-activity-id-3",
    userId: "sportsfan",
    userName: "Carlos Rodriguez",
    userEmail: "sports.fan@test.com",
    date: new Date("2025-12-17T06:00:00Z"),
    participants: 2,
    status: "confirmed",
    totalPrice: 600.0,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    activityId: activityIds[3] || "sample-activity-id-4",
    userId: "yogalover",
    userName: "Maria Garcia",
    userEmail: "yoga.lover@test.com",
    date: new Date("2025-12-18T19:00:00Z"),
    participants: 4,
    status: "pending",
    totalPrice: 1600.0,
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    activityId: activityIds[4] || "sample-activity-id-5",
    userId: "runner",
    userName: "Pedro Martinez",
    userEmail: "runner@test.com",
    date: new Date("2025-12-20T08:00:00Z"),
    participants: 1,
    status: "confirmed",
    totalPrice: 450.0,
    createdAt: new Date(),
    updatedAt: new Date()
  }
];

const reservationResult = db.reservas.insertMany(reservations);
print(`✓ Insertadas ${reservationResult.insertedIds.length} reservas en reservasdb`);

// Estadísticas de reservas
print("\n→ Estadísticas de reservas:");
const stats = db.reservas.aggregate([
  {
    $group: {
      _id: "$status",
      count: { $sum: 1 },
      totalRevenue: { $sum: "$totalPrice" }
    }
  }
]).toArray();

stats.forEach(stat => {
  print(`  - ${stat._id}: ${stat.count} reservas, $${stat.totalRevenue} total`);
});

print("\n===========================================");
print("✅ Seeding completado exitosamente!");
print("===========================================\n");

print("📊 Resumen:");
print(`  - Actividades: ${activityIds.length} insertadas`);
print(`  - Reservas: ${reservations.length} insertadas`);
print("\nPuedes verificar con:");
print("  mongosh activitiesdb --eval 'db.activities.countDocuments()'");
print("  mongosh reservasdb --eval 'db.reservas.countDocuments()'");
