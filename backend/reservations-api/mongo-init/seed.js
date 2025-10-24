// Seed some data if the collection is empty
db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE || "demo");
const col = db.getCollection("reservas");

if (col.countDocuments() === 0) {
  col.insertMany([
    {
      user_id: 1,
      actividad: "Visita guiada al museo",
      date: ISODate("2026-01-05T15:00:00Z"),
      status: "cancelada",
      created_at: new Date(),
      updated_at: new Date(),
    },
    {
      user_id: 2,
      actividad: "Clase de cocina italiana",
      date: ISODate("2026-02-10T18:30:00Z"),
      status: "confirmada",
      created_at: new Date(),
      updated_at: new Date(),
    },
    {
      user_id: 3,
      actividad: "Excursión al parque nacional",
      date: ISODate("2026-03-22T09:00:00Z"),
      status: "pendiente",
      created_at: new Date(),
      updated_at: new Date(),
    },
    {
      user_id: 4,
      actividad: "Taller de fotografía urbana",
      date: ISODate("2026-04-15T14:00:00Z"),
      status: "confirmada",
      created_at: new Date(),
      updated_at: new Date(),
    },
    {
      user_id: 5,
      actividad: "Clase de yoga al aire libre",
      date: ISODate("2026-05-01T08:00:00Z"),
      status: "pendiente",
      created_at: new Date(),
      updated_at: new Date(),
    },
  ]);
  print("Seeded initial reservas");
} else {
  print("Reservas already present, skipping seed");
}
