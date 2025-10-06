import React, { useState, useEffect } from 'react';
import {
  Box,
  TextField,
  Typography,
  Chip,
  FormHelperText,
  InputAdornment,
  Tooltip,
} from '@mui/material';
import {
  DatePicker,
  LocalizationProvider,
} from '@mui/x-date-pickers';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import {
  Cake as CakeIcon,
  Today as TodayIcon,
  Info as InfoIcon,
} from '@mui/icons-material';
import dayjs, { Dayjs } from 'dayjs';

interface BirthDatePickerProps {
  value?: string | null;
  onChange: (date: string | null) => void;
  disabled?: boolean;
  error?: boolean;
  helperText?: string;
  label?: string;
  fullWidth?: boolean;
}

interface AgeInfo {
  years: number;
  months: number;
  days: number;
  totalDays: number;
  zodiacSign: string;
  isValid: boolean;
}

// Zodiac signs calculation
const getZodiacSign = (month: number, day: number): string => {
  const zodiacSigns = [
    { sign: 'Capricorn ♑', start: [12, 22], end: [1, 19] },
    { sign: 'Aquarius ♒', start: [1, 20], end: [2, 18] },
    { sign: 'Pisces ♓', start: [2, 19], end: [3, 20] },
    { sign: 'Aries ♈', start: [3, 21], end: [4, 19] },
    { sign: 'Taurus ♉', start: [4, 20], end: [5, 20] },
    { sign: 'Gemini ♊', start: [5, 21], end: [6, 20] },
    { sign: 'Cancer ♋', start: [6, 21], end: [7, 22] },
    { sign: 'Leo ♌', start: [7, 23], end: [8, 22] },
    { sign: 'Virgo ♍', start: [8, 23], end: [9, 22] },
    { sign: 'Libra ♎', start: [9, 23], end: [10, 22] },
    { sign: 'Scorpio ♏', start: [10, 23], end: [11, 21] },
    { sign: 'Sagittarius ♐', start: [11, 22], end: [12, 21] },
  ];

  for (const zodiac of zodiacSigns) {
    const [startMonth, startDay] = zodiac.start;
    const [endMonth, endDay] = zodiac.end;
    
    if (
      (month === startMonth && day >= startDay) ||
      (month === endMonth && day <= endDay) ||
      (startMonth === 12 && month === 1 && day <= endDay) // Handle year boundary
    ) {
      return zodiac.sign;
    }
  }
  
  return 'Unknown';
};

const calculateAge = (birthDate: Dayjs): AgeInfo => {
  const today = dayjs();
  const birth = birthDate;
  
  if (birth.isAfter(today)) {
    return {
      years: 0,
      months: 0,
      days: 0,
      totalDays: 0,
      zodiacSign: 'Invalid',
      isValid: false,
    };
  }

  const years = today.diff(birth, 'years');
  const months = today.diff(birth.add(years, 'years'), 'months');
  const days = today.diff(birth.add(years, 'years').add(months, 'months'), 'days');
  const totalDays = today.diff(birth, 'days');
  
  const zodiacSign = getZodiacSign(birth.month() + 1, birth.date());

  return {
    years,
    months,
    days,
    totalDays,
    zodiacSign,
    isValid: years >= 0 && years <= 150,
  };
};

const BirthDatePicker: React.FC<BirthDatePickerProps> = ({
  value,
  onChange,
  disabled = false,
  error = false,
  helperText,
  label = "Birth Date",
  fullWidth = true,
}) => {
  const [selectedDate, setSelectedDate] = useState<Dayjs | null>(
    value ? dayjs(value) : null
  );
  const [ageInfo, setAgeInfo] = useState<AgeInfo | null>(null);
  const [internalError, setInternalError] = useState<string>('');

  useEffect(() => {
    if (selectedDate) {
      const age = calculateAge(selectedDate);
      setAgeInfo(age);
      
      // Validation
      if (!age.isValid) {
        setInternalError('Please enter a valid birth date');
      } else if (age.years < 13) {
        setInternalError('You must be at least 13 years old');
      } else if (age.years > 150) {
        setInternalError('Please enter a realistic birth date');
      } else {
        setInternalError('');
      }
    } else {
      setAgeInfo(null);
      setInternalError('');
    }
  }, [selectedDate]);

  const handleDateChange = (newDate: Dayjs | null) => {
    setSelectedDate(newDate);
    onChange(newDate ? newDate.format('YYYY-MM-DD') : null);
  };

  const getNextBirthday = (): string => {
    if (!selectedDate) return '';
    
    const today = dayjs();
    const thisYearBirthday = selectedDate.year(today.year());
    const nextBirthday = thisYearBirthday.isBefore(today) 
      ? thisYearBirthday.add(1, 'year') 
      : thisYearBirthday;
    
    const daysUntil = nextBirthday.diff(today, 'days');
    
    if (daysUntil === 0) return 'Today! 🎉';
    if (daysUntil === 1) return 'Tomorrow! 🎂';
    if (daysUntil <= 30) return `In ${daysUntil} days 🎈`;
    
    return nextBirthday.format('MMM DD');
  };

  const maxDate = dayjs().subtract(13, 'years');
  const minDate = dayjs().subtract(150, 'years');

  return (
    <LocalizationProvider dateAdapter={AdapterDayjs}>
      <Box>
        <DatePicker
          label={label}
          value={selectedDate}
          onChange={handleDateChange}
          disabled={disabled}
          maxDate={maxDate}
          minDate={minDate}
          format="MMM DD, YYYY"
          slots={{
            textField: (params) => (
              <TextField
                {...params}
                fullWidth={fullWidth}
                error={error || !!internalError}
                helperText={helperText || internalError}
                InputProps={{
                  ...params.InputProps,
                  startAdornment: (
                    <InputAdornment position="start">
                      <CakeIcon color={error || internalError ? 'error' : 'action'} />
                    </InputAdornment>
                  ),
                }}
              />
            ),
          }}
          slotProps={{
            popper: {
              placement: 'bottom-start',
            },
            desktopPaper: {
              sx: {
                '& .MuiPickersCalendarHeader-root': {
                  backgroundColor: 'primary.main',
                  color: 'primary.contrastText',
                },
              },
            },
          }}
        />

        {/* Age Information Display */}
        {ageInfo && ageInfo.isValid && !internalError && (
          <Box sx={{ mt: 2, display: 'flex', flexWrap: 'wrap', gap: 1, alignItems: 'center' }}>
            <Chip
              icon={<TodayIcon />}
              label={`${ageInfo.years} years old`}
              color="primary"
              variant="outlined"
              size="small"
            />
            
            {ageInfo.zodiacSign !== 'Unknown' && (
              <Chip
                label={ageInfo.zodiacSign}
                color="secondary"
                variant="outlined"
                size="small"
              />
            )}
            
            <Tooltip 
              title={`Next birthday: ${getNextBirthday()}`}
              arrow
            >
              <Chip
                icon={<InfoIcon />}
                label={`${ageInfo.totalDays.toLocaleString()} days lived`}
                variant="outlined"
                size="small"
                sx={{ cursor: 'help' }}
              />
            </Tooltip>
          </Box>
        )}

        {/* Detailed Age Breakdown */}
        {ageInfo && ageInfo.isValid && !internalError && (
          <Box sx={{ mt: 1 }}>
            <Typography variant="caption" color="text.secondary">
              Exact age: {ageInfo.years} years, {ageInfo.months} months, {ageInfo.days} days
              {getNextBirthday() && ` • Next birthday: ${getNextBirthday()}`}
            </Typography>
          </Box>
        )}

        {/* Additional Helper Text */}
        {!selectedDate && !disabled && (
          <FormHelperText sx={{ mt: 1, color: 'text.secondary' }}>
            Click the calendar icon to select your birth date. This helps us personalize your experience.
          </FormHelperText>
        )}
      </Box>
    </LocalizationProvider>
  );
};

export default BirthDatePicker;
