import dayjs from "dayjs";

const parseDate = (
  date?: string | number | Date | dayjs.Dayjs | null | undefined,
) => {
  return date ? dayjs(date) : null;
};

const formatDate = (
  date?: string | number | Date | dayjs.Dayjs | null | undefined,
) => {
  if (!date) {
    return "null";
  }
  return parseDate(date)?.format("YYYY-MM-DD") || "null";
};

const timestampForFilename = (date: Date | number) => {
  const nowParsed = dayjs(date);
  return nowParsed.format("YYYYMMDDTHHmmss");
};

export { formatDate, parseDate, timestampForFilename };
