-- AlterTable
ALTER TABLE "Event" ALTER COLUMN "createdAt" SET DEFAULT CLOCK_TIMESTAMP(),
ALTER COLUMN "createdAt" SET DATA TYPE TIMESTAMP(6);