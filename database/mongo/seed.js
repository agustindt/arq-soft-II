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

// Actividades de prueba variadas (usando nombres de campos correctos para el modelo)
const activities = [
  {
    name: "Fútbol 5 Nocturno",
    description: "Partido de fútbol 5 en cancha sintética bajo techo. Ideal para después del trabajo. Incluye arbitraje y balón.",
    category: "football",
    difficulty: "intermediate",
    duration: 90,
    max_capacity: 10,
    location: "Polideportivo Central, Cancha 3",
    price: 500.0,
    image_url: "https://images.unsplash.com/photo-1431324155629-1a6deb1dec8d?w=800",
    instructor: "Carlos Mendoza",
    schedule: ["Lunes 20:00", "Miércoles 20:00", "Viernes 20:00"],
    equipment: ["Balón", "Cancha", "Arbitraje"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Yoga Matutino en el Parque",
    description: "Sesión de yoga al aire libre para principiantes y avanzados. Trae tu mat. Instructor certificado.",
    category: "yoga",
    difficulty: "beginner",
    duration: 60,
    max_capacity: 20,
    location: "Parque de la Reserva, Zona Verde",
    price: 250.0,
    image_url: "https://images.unsplash.com/photo-1506126613408-eca07ce68773?w=800",
    instructor: "Ana García",
    schedule: ["Martes 07:00", "Jueves 07:00", "Sábado 08:00"],
    equipment: ["Mat de yoga"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Running Club 10K",
    description: "Entrenamiento grupal de running para preparación de 10K. Incluye plan de entrenamiento y análisis de técnica.",
    category: "running",
    difficulty: "advanced",
    duration: 120,
    max_capacity: 30,
    location: "Malecón Costa Verde",
    price: 300.0,
    image_url: "https://images.unsplash.com/photo-1571008887538-b36bb32f4571?w=800",
    instructor: "Marco Silva",
    schedule: ["Lunes 06:00", "Miércoles 06:00", "Sábado 07:00"],
    equipment: ["Zapatillas de running"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Básquetbol 3x3",
    description: "Torneo de básquetbol 3x3 estilo streetball. Premios para los ganadores.",
    category: "basketball",
    difficulty: "intermediate",
    duration: 120,
    max_capacity: 12,
    location: "Cancha Municipal Norte",
    price: 400.0,
    image_url: "https://images.unsplash.com/photo-1519861531473-9200262198bf?w=800",
    instructor: "Luis Rodríguez",
    schedule: ["Viernes 18:00", "Domingo 16:00"],
    equipment: ["Balón", "Cancha"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Natación para Principiantes",
    description: "Clases de natación para adultos que están empezando. Incluye instructor y equipo básico.",
    category: "swimming",
    difficulty: "beginner",
    duration: 60,
    max_capacity: 8,
    location: "Club Náutico, Piscina Olímpica",
    price: 450.0,
    image_url: "https://images.unsplash.com/photo-1530549387789-4c1017266635?w=800",
    instructor: "Patricia López",
    schedule: ["Lunes 19:00", "Miércoles 19:00", "Viernes 19:00"],
    equipment: ["Gorro", "Lentes", "Tabla"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Ciclismo de Montaña",
    description: "Ruta de mountain bike por cerros locales. Nivel intermedio-avanzado. Trae tu propia bici.",
    category: "cycling",
    difficulty: "advanced",
    duration: 180,
    max_capacity: 15,
    location: "Inicio: Parque Zonal Huáscar",
    price: 350.0,
    image_url: "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800",
    instructor: "Roberto Martínez",
    schedule: ["Sábado 08:00", "Domingo 08:00"],
    equipment: ["Bicicleta de montaña", "Casco", "Agua"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Tenis Singles - Torneo Amateur",
    description: "Torneo de tenis singles amateur. Canchas de arcilla. Sistema de eliminación directa.",
    category: "tennis",
    difficulty: "intermediate",
    duration: 90,
    max_capacity: 16,
    location: "Club de Tenis Los Alamos",
    price: 600.0,
    image_url: "https://images.unsplash.com/photo-1622279457486-62dcc4a431f3?w=800",
    instructor: "Fernando Torres",
    schedule: ["Sábado 10:00", "Domingo 10:00"],
    equipment: ["Raqueta", "Pelotas", "Cancha"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Volleyball Playa Mixto",
    description: "Partidos de volleyball en la playa. Equipos mixtos. Ambiente relajado y divertido.",
    category: "volleyball",
    difficulty: "beginner",
    duration: 120,
    max_capacity: 12,
    location: "Playa Agua Dulce",
    price: 200.0,
    image_url: "https://images.unsplash.com/photo-1612872087720-bb876e2e67d1?w=800",
    instructor: "Sofía Ramírez",
    schedule: ["Sábado 16:00", "Domingo 16:00"],
    equipment: ["Red", "Balón"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "CrossFit Extremo",
    description: "WOD intensivo de CrossFit. Solo para atletas avanzados. Incluye coach certificado.",
    category: "fitness",
    difficulty: "advanced",
    duration: 90,
    max_capacity: 15,
    location: "Box CrossFit Revolution",
    price: 550.0,
    image_url: "https://images.unsplash.com/photo-1534438327276-14e5300c3a48?w=800",
    instructor: "Diego Herrera",
    schedule: ["Lunes 19:00", "Miércoles 19:00", "Viernes 19:00"],
    equipment: ["Pesas", "Kettlebells", "Barras"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Paddle Tennis Dobles",
    description: "Sesión de pádel en modalidad dobles. Todos los niveles bienvenidos. Raquetas disponibles.",
    category: "paddle",
    difficulty: "beginner",
    duration: 90,
    max_capacity: 8,
    location: "Centro Deportivo Las Palmas",
    price: 400.0,
    image_url: "https://images.unsplash.com/photo-1622163642999-8bd0876ad1f3?w=800",
    instructor: "Carmen Vásquez",
    schedule: ["Martes 20:00", "Jueves 20:00"],
    equipment: ["Raquetas", "Pelotas", "Cancha"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Fútbol 11 Profesional",
    description: "Partido de fútbol 11 en cancha profesional. Nivel avanzado. Incluye arbitraje profesional.",
    category: "football",
    difficulty: "advanced",
    duration: 90,
    max_capacity: 22,
    location: "Estadio Municipal",
    price: 800.0,
    image_url: "https://images.unsplash.com/photo-1574629810360-7efbbe195018?w=800",
    instructor: "Miguel Ángel",
    schedule: ["Domingo 10:00"],
    equipment: ["Balón", "Cancha", "Arbitraje"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Yoga Restaurativo",
    description: "Sesión de yoga restaurativo para relajación y recuperación. Perfecto para después del entrenamiento.",
    category: "yoga",
    difficulty: "beginner",
    duration: 75,
    max_capacity: 15,
    location: "Estudio Zen Yoga",
    price: 300.0,
    image_url: "https://images.unsplash.com/photo-1506126613408-eca07ce68773?w=800",
    instructor: "Laura Fernández",
    schedule: ["Lunes 20:00", "Miércoles 20:00"],
    equipment: ["Mat", "Almohadas", "Mantas"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Maratón de Entrenamiento",
    description: "Entrenamiento intensivo para maratón. Plan de 16 semanas. Incluye seguimiento personalizado.",
    category: "running",
    difficulty: "advanced",
    duration: 150,
    max_capacity: 25,
    location: "Parque Metropolitano",
    price: 450.0,
    image_url: "https://images.unsplash.com/photo-1571008887538-b36bb32f4571?w=800",
    instructor: "Eduardo Castro",
    schedule: ["Sábado 06:00", "Domingo 06:00"],
    equipment: ["Zapatillas", "Reloj GPS"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Básquetbol Femenino",
    description: "Liga de básquetbol exclusiva para mujeres. Todos los niveles. Ambiente inclusivo y motivador.",
    category: "basketball",
    difficulty: "beginner",
    duration: 90,
    max_capacity: 20,
    location: "Polideportivo Sur",
    price: 350.0,
    image_url: "https://images.unsplash.com/photo-1519861531473-9200262198bf?w=800",
    instructor: "Valentina Morales",
    schedule: ["Martes 19:00", "Jueves 19:00"],
    equipment: ["Balón", "Cancha"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Natación Competitiva",
    description: "Entrenamiento de natación competitiva. Técnica avanzada y preparación para competencias.",
    category: "swimming",
    difficulty: "advanced",
    duration: 90,
    max_capacity: 12,
    location: "Centro Acuático Olímpico",
    price: 600.0,
    image_url: "https://images.unsplash.com/photo-1530549387789-4c1017266635?w=800",
    instructor: "Ricardo Paredes",
    schedule: ["Lunes 06:00", "Miércoles 06:00", "Viernes 06:00"],
    equipment: ["Gorro", "Lentes", "Aletas"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Ciclismo Urbano",
    description: "Rutas de ciclismo por la ciudad. Nivel principiante. Ideal para conocer la ciudad mientras haces ejercicio.",
    category: "cycling",
    difficulty: "beginner",
    duration: 120,
    max_capacity: 20,
    location: "Plaza Central",
    price: 250.0,
    image_url: "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800",
    instructor: "Andrés Gutiérrez",
    schedule: ["Sábado 09:00", "Domingo 09:00"],
    equipment: ["Bicicleta", "Casco"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Tenis Dobles Social",
    description: "Sesión de tenis dobles para jugadores sociales. Ambiente amigable y relajado.",
    category: "tennis",
    difficulty: "beginner",
    duration: 90,
    max_capacity: 8,
    location: "Club Social Deportivo",
    price: 400.0,
    image_url: "https://images.unsplash.com/photo-1622279457486-62dcc4a431f3?w=800",
    instructor: "Gloria Sánchez",
    schedule: ["Sábado 11:00", "Domingo 11:00"],
    equipment: ["Raquetas", "Pelotas"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Volleyball Competitivo",
    description: "Liga de volleyball competitivo. Nivel intermedio-avanzado. Torneos mensuales.",
    category: "volleyball",
    difficulty: "intermediate",
    duration: 120,
    max_capacity: 24,
    location: "Gimnasio Municipal",
    price: 450.0,
    image_url: "https://images.unsplash.com/photo-1612872087720-bb876e2e67d1?w=800",
    instructor: "Jorge Méndez",
    schedule: ["Martes 20:00", "Jueves 20:00"],
    equipment: ["Red", "Balón", "Cancha"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "HIIT Training",
    description: "Entrenamiento de alta intensidad (HIIT). Quema calorías y mejora tu condición física en 30 minutos.",
    category: "fitness",
    difficulty: "intermediate",
    duration: 45,
    max_capacity: 20,
    location: "Gimnasio FitZone",
    price: 350.0,
    image_url: "https://images.unsplash.com/photo-1534438327276-14e5300c3a48?w=800",
    instructor: "Natalia Ríos",
    schedule: ["Lunes 18:00", "Miércoles 18:00", "Viernes 18:00"],
    equipment: ["Pesas", "Colchonetas"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Paddle Tennis Individual",
    description: "Clases individuales de pádel. Enfoque en técnica y estrategia. Todos los niveles.",
    category: "paddle",
    difficulty: "intermediate",
    duration: 60,
    max_capacity: 4,
    location: "Club Pádel Elite",
    price: 500.0,
    image_url: "https://images.unsplash.com/photo-1622163642999-8bd0876ad1f3?w=800",
    instructor: "Alejandro Vega",
    schedule: ["Lunes 19:00", "Miércoles 19:00"],
    equipment: ["Raquetas", "Pelotas"],
    is_active: true,
    created_by: 1,
    created_at: new Date(),
    updated_at: new Date()
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
