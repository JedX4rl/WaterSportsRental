CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS people
(
    ID               SERIAL PRIMARY KEY,
    FirstName        VARCHAR(255),
    LastName         VARCHAR(255),
    Address          VARCHAR(255),
    PhoneNumber      VARCHAR(20),
    RegistrationDate TIMESTAMP,
    Role             VARCHAR(255),
    Email            VARCHAR(255) UNIQUE NOT NULL,
    Password         VARCHAR(255)        NOT NULL
);

CREATE TABLE IF NOT EXISTS dumps
(
    ID       SERIAL PRIMARY KEY,
    filePath VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS locations
(
    LocationID  SERIAL PRIMARY KEY,
    Country     VARCHAR(255) NOT NULL,
    City        VARCHAR(255) NOT NULL,
    Address     VARCHAR(255) NOT NULL,
    OpeningTime TIME         NOT NULL,
    ClosingTime TIME         NOT NULL,
    PhoneNumber VARCHAR(50),
    CONSTRAINT unique_locations UNIQUE (Country, City, Address, PhoneNumber)
);

CREATE TABLE IF NOT EXISTS inventory
(
    InventoryID       SERIAL PRIMARY KEY,
    Type              VARCHAR(255) NOT NULL,
    Brand             VARCHAR(255) NOT NULL,
    Model             VARCHAR(255) NOT NULL,
    Year              INT CHECK ( Year >= 1990 AND Year <= EXTRACT(YEAR FROM CURRENT_DATE)),
    RentalPricePerDay DECIMAL      NOT NULL,
    ImageURL          TEXT,
    CONSTRAINT unique_inventory UNIQUE (Type, Brand, Model, Year)
);


CREATE TABLE IF NOT EXISTS reviews
(
    ReviewID   SERIAL PRIMARY KEY,
    Name       VARCHAR(255) NOT NULL,
    Rating     INT CHECK ( Rating >= 1 AND Rating <= 5 ),
    Comment    TEXT,
    ReviewDate DATE         NOT NULL
);


CREATE TABLE IF NOT EXISTS payments
(
    PaymentID     SERIAL PRIMARY KEY,
    Amount        DECIMAL     NOT NULL,
    PaymentDate   DATE        NOT NULL,
    PaymentMethod VARCHAR(40) NOT NULL,
    Status        BOOLEAN     NOT NULL
);


CREATE TABLE IF NOT EXISTS rentalAgreements
(
    RentalAgreementID SERIAL PRIMARY KEY,
    PaymentID         INT,
    CustomerId        INT     NOT NULL,
    DateOfAgreement   DATE    NOT NULL,
    Status            BOOLEAN NOT NULL,
    FOREIGN KEY (CustomerId) references people (ID),
    FOREIGN KEY (PaymentID) references payments (PaymentID)
);

CREATE TABLE IF NOT EXISTS rentedEquipment
(
    RentedID          SERIAL PRIMARY KEY,
    RentalAgreementId INT  NOT NULL,
    ReviewID          INT,
    InventoryID       INT  NOT NULL,
    StartDate         DATE NOT NULL,
    EndDate           DATE NOT NULL,
    FOREIGN KEY (RentalAgreementId) references rentalAgreements (RentalAgreementID),
    FOREIGN KEY (InventoryID) references inventory (InventoryID),
    FOREIGN KEY (ReviewID) references reviews (ReviewID)
);


CREATE TABLE IF NOT EXISTS rented_equipment_log
(
    id           SERIAL PRIMARY KEY,
    inventory_id INT NOT NULL,
    start_date   DATE       NOT NULL,
    end_date     DATE       NOT NULL
);


CREATE OR REPLACE FUNCTION log_rented_equipment()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO rented_equipment_log (inventory_id, start_date, end_date)
    VALUES (NEW.InventoryID, NEW.StartDate, NEW.EndDate);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER after_insert_rented_equipment
    AFTER INSERT ON rentedEquipment
    FOR EACH ROW
EXECUTE FUNCTION log_rented_equipment();


CREATE OR REPLACE FUNCTION get_all_rented_equipment_logs()
    RETURNS TABLE(
                     id INT,
                     inventory_id INT,
                     start_date DATE,
                     end_date DATE
                 ) AS $$
BEGIN
    RETURN QUERY
        SELECT rented_equipment_log.id, rented_equipment_log.inventory_id, rented_equipment_log.start_date, rented_equipment_log.end_date
        FROM rented_equipment_log;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_rented_items(p_user_id INT)
    RETURNS TABLE
            (
                id          INT,
                type        VARCHAR(255),
                brand       VARCHAR(255),
                model       VARCHAR(255),
                year        INT,
                total_price DECIMAL,
                image       TEXT,
                start_date  DATE,
                end_date    DATE
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT i.InventoryID                                     AS id,
               i.Type                                            AS type,
               i.Brand                                           AS brand,
               i.Model                                           AS model,
               i.Year                                            AS year,
               i.RentalPricePerDay * (re.EndDate - re.StartDate) AS total_price,
               i.ImageURL                                        AS image,
               re.StartDate                                      AS start_date,
               re.EndDate                                        AS end_date
        FROM rentedEquipment re
                 JOIN rentalAgreements ra ON re.RentalAgreementId = ra.RentalAgreementID
                 JOIN inventory i ON re.InventoryID = i.InventoryID
        WHERE ra.CustomerId = p_user_id;
END;
$$ LANGUAGE plpgsql;



    CREATE TABLE IF NOT EXISTS productAvailability
    (
        Id          SERIAL PRIMARY KEY,
        InventoryID INT NOT NULL,
        LocationId  INT NOT NULL,
        Number      INT NOT NULL,
        FOREIGN KEY (InventoryID) references inventory (InventoryID),
        FOREIGN KEY (LocationId) references locations (locationid)
    );

CREATE OR REPLACE FUNCTION insert_dump(p_filePath VARCHAR)
    RETURNS VOID AS
$$
BEGIN
    INSERT INTO dumps (filePath)
    VALUES (p_filePath)
    ON CONFLICT (filePath) DO NOTHING;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_all_dumps()
    RETURNS TABLE
            (
                filePath VARCHAR
            )
AS
$$
BEGIN
    RETURN QUERY SELECT dumps.filePath FROM dumps;
END;
$$ LANGUAGE plpgsql;


--RENTAL AGREEMENT--

CREATE OR REPLACE FUNCTION create_rental_agreement(
    p_customer_id INT,
    p_date_of_agreement TIMESTAMP,
    p_status BOOLEAN
) RETURNS INT AS
$$
DECLARE
    v_rental_agreement_id INT;
BEGIN
    INSERT INTO rentalagreements (customerid, dateofagreement, status)
    VALUES (p_customer_id, p_date_of_agreement, p_status)
    RETURNING rentalagreementid INTO v_rental_agreement_id;
    RETURN v_rental_agreement_id;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION delete_rental_agreement(p_agreement_id INT) RETURNS VOID AS
$$
BEGIN
    DELETE FROM rentalagreements WHERE rentalagreementid = p_agreement_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_all_rental_agreements(p_user_id INT)
    RETURNS TABLE
            (
                dateofagreement DATE,
                equipmentid     INT,
                startdate       DATE,
                enddate         DATE,
                status          BOOLEAN
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT ra.dateofagreement, re.inventoryid, re.startdate, re.enddate, ra.status
        FROM rentalagreements ra
                 LEFT JOIN rentedequipment re ON ra.rentalagreementid = re.rentalagreementid
        WHERE ra.customerid = p_user_id;
END;
$$ LANGUAGE plpgsql;


--RENTAL AGREEMENT--

--AUTH--

CREATE OR REPLACE FUNCTION create_user(
    p_email TEXT,
    p_password TEXT,
    p_registrationdate TIMESTAMP
) RETURNS INT AS
$$
DECLARE
    new_user_id INT;
BEGIN
    INSERT INTO people (email, password, registrationdate, role)
    VALUES (p_email, p_password, p_registrationdate, 'user')
    RETURNING id INTO new_user_id;

    RETURN new_user_id;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_user_by_email(user_email TEXT)
    RETURNS TABLE
            (
                p_id               INT,
                p_first_name       VARCHAR(255),
                p_last_name        VARCHAR(255),
                p_address          VARCHAR(255),
                p_phoneNumber      VARCHAR(20),
                p_registrationDate TIMESTAMP,
                p_role             VARCHAR(255),
                p_email            VARCHAR(255),
                p_password         VARCHAR(255)
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT id               AS p_id,
               firstname        AS p_first_name,
               lastname         AS p_last_name,
               address          AS p_address,
               phonenumber      AS p_phoneNumber,
               registrationdate AS p_registrationDate,
               role             AS p_role,
               email            AS p_email,
               password         AS p_password
        FROM people
        WHERE Email = user_email;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_user_by_id(user_id INT)
    RETURNS TABLE
            (
                p_id               INT,
                p_first_name       VARCHAR(255),
                p_last_name        VARCHAR(255),
                p_address          VARCHAR(255),
                p_phoneNumber      VARCHAR(20),
                p_registrationDate TIMESTAMP,
                p_role             VARCHAR(255),
                p_email            VARCHAR(255),
                p_password         VARCHAR(255)
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT id               AS p_id,
               firstname        AS p_first_name,
               lastname         AS p_last_name,
               address          AS p_address,
               phonenumber      AS p_phoneNumber,
               registrationdate AS p_registrationDate,
               role             AS p_role,
               email            AS p_email,
               password         AS p_password
        FROM people
        WHERE ID = user_id;
END
$$ LANGUAGE plpgsql;


--AUTH--


--USER--

CREATE OR REPLACE FUNCTION update_user(
    p_id INT,
    p_first_name VARCHAR(255) DEFAULT NULL,
    p_last_name VARCHAR(255) DEFAULT NULL,
    p_address VARCHAR(255) DEFAULT NULL,
    p_phone_number VARCHAR(20) DEFAULT NULL,
    p_role VARCHAR(255) DEFAULT NULL
) RETURNS VOID AS
$$
BEGIN
    UPDATE people
    SET FirstName   = COALESCE(p_first_name, FirstName),
        LastName    = COALESCE(p_last_name, LastName),
        Address     = COALESCE(p_address, Address),
        PhoneNumber = COALESCE(p_phone_number, PhoneNumber),
        Role        = COALESCE(p_role, Role)
    WHERE ID = p_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION give_role(
    p_id INT,
    p_role VARCHAR(255) DEFAULT NULL
) RETURNS VOID AS
$$
BEGIN
    UPDATE people
    SET Role = COALESCE(p_role, Role)
    WHERE ID = p_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_user(user_id INT)
    RETURNS TABLE
            (
                p_id                INT,
                p_email             VARCHAR(255),
                p_first_name        VARCHAR(255),
                p_last_name         VARCHAR(255),
                p_address           VARCHAR(255),
                p_phone_number      VARCHAR(20),
                p_registration_date TIMESTAMP,
                p_role              VARCHAR(255)
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT ID,
               Email,
               FirstName,
               LastName,
               Address,
               PhoneNumber,
               RegistrationDate,
               Role
        FROM people
        WHERE ID = user_id;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_all_users()
    RETURNS TABLE
            (
                p_id                INT,
                p_email             VARCHAR(255),
                p_first_name        VARCHAR(255),
                p_last_name         VARCHAR(255),
                p_address           VARCHAR(255),
                p_phone_number      VARCHAR(20),
                p_registration_date TIMESTAMP,
                p_role              VARCHAR(255)
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT ID               AS p_id,
               Email               p_email,
               FirstName        AS p_first_name,
               LastName         AS p_last_name,
               Address          AS p_address,
               PhoneNumber      AS p_phone_number,
               RegistrationDate AS p_registration_date,
               Role             AS p_role
        FROM people
        ORDER BY ID;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION delete_user(p_id INT)
    RETURNS VOID AS
$$
BEGIN
    DELETE FROM people WHERE ID = p_id;
END;
$$ LANGUAGE plpgsql;


--USER--


CREATE OR REPLACE FUNCTION add_item(
    p_type VARCHAR(255),
    p_brand VARCHAR(255),
    p_model VARCHAR(255),
    p_year INT,
    p_price_per_day NUMERIC,
    p_image VARCHAR(255)
) RETURNS INT AS
$$
DECLARE
    new_item_id INT;
BEGIN
    INSERT INTO inventory(type, brand, model, year, rentalpriceperday, imageUrl)
    VALUES (p_type, p_brand, p_model, p_year, p_price_per_day, p_image)
    RETURNING InventoryID INTO new_item_id;
    RETURN new_item_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION delete_item(p_id INT)
    RETURNS VOID AS
$$
BEGIN
    DELETE FROM inventory WHERE InventoryID = p_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_item(p_id INT)
    RETURNS TABLE
            (
                inventory_id         INT,
                type                 TEXT,
                brand                TEXT,
                model                TEXT,
                year                 INT,
                rental_price_per_day FLOAT,
                image_url            TEXT
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT inventoryid, type, brand, model, year, rentalpriceperday, imageurl
        FROM inventory
        WHERE inventoryid = p_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_filtered_items(
    p_filters JSONB DEFAULT '{}'
)
    RETURNS TABLE
            (
                inventory_id         INT,
                type                 VARCHAR(255),
                brand                VARCHAR(255),
                model                VARCHAR(255),
                year                 INT,
                rental_price_per_day NUMERIC,
                image_url            TEXT
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT DISTINCT i.inventoryid,
                        i.type,
                        i.brand,
                        i.model,
                        i.year,
                        i.rentalpriceperday,
                        i.imageurl
        FROM inventory i
                 LEFT JOIN productavailability pa ON i.inventoryid = pa.inventoryid
        WHERE (
            p_filters -> 'locationId' IS NULL OR
            pa.locationid = ANY (ARRAY(
                    SELECT jsonb_array_elements_text(p_filters -> 'locationId')::INT
                                 ))
            )
          AND (
            p_filters -> 'type' IS NULL OR
            i.type = ANY (ARRAY(SELECT jsonb_array_elements_text(p_filters -> 'type')))
            )
          AND (
            p_filters -> 'brand' IS NULL OR
            i.brand = ANY (ARRAY(SELECT jsonb_array_elements_text(p_filters -> 'brand')))
            )
          AND (
            p_filters -> 'model' IS NULL OR
            i.model = ANY (ARRAY(SELECT jsonb_array_elements_text(p_filters -> 'model')))
            )
          AND (
            p_filters -> 'year' IS NULL OR
            i.year = ANY (ARRAY(SELECT jsonb_array_elements_text(p_filters -> 'year')::INT))
            )
          AND (
            p_filters -> 'priceMin' IS NULL OR
            i.rentalpriceperday >= (p_filters ->> 'priceMin')::FLOAT
            )
          AND (
            p_filters -> 'priceMax' IS NULL OR
            i.rentalpriceperday <= (p_filters ->> 'priceMax')::FLOAT
            );
END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_all_items(p_filters JSONB DEFAULT '{}')
    RETURNS TABLE
            (
                inventory_id         INT,
                type                 VARCHAR(255),
                brand                VARCHAR(255),
                model                VARCHAR(255),
                year                 INT,
                rental_price_per_day NUMERIC,
                image_url            TEXT
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT DISTINCT i.inventoryid,
                        i.type,
                        i.brand,
                        i.model,
                        i.year,
                        i.rentalpriceperday,
                        i.imageurl
        FROM inventory i
                 LEFT JOIN productavailability pa ON i.inventoryid = pa.inventoryid
        ORDER BY InventoryID;
END ;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_available_products(p_inventory_id INT)
    RETURNS TABLE
            (
                product_id  INT,
                location_id INT,
                location    VARCHAR(255),
                number      INT
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT pa.inventoryid AS product_id,
               l.locationid   AS location_id,
               l.address      AS location,
               pa.number
        FROM productavailability pa
                 RIGHT JOIN public.locations l ON pa.locationid = l.locationid
        WHERE pa.inventoryid IS NOT NULL
          AND pa.inventoryid = p_inventory_id;
END
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION update_item(
    p_id INT,
    p_type VARCHAR(255),
    p_brand VARCHAR(255),
    p_model VARCHAR(255),
    p_year INT,
    p_price_per_day NUMERIC,
    p_image VARCHAR(255)
) RETURNS VOID AS
$$
BEGIN
    UPDATE inventory
    SET Type              = COALESCE(p_type, Type),
        Brand             = COALESCE(p_brand, Brand),
        Model             = COALESCE(p_model, Model),
        Year              = COALESCE(p_year, Year),
        RentalPricePerDay = COALESCE(p_price_per_day, RentalPricePerDay),
        ImageURL          = COALESCE(p_image, ImageURL)
    WHERE InventoryID = p_id;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION rent_equipment(
    p_rental_agreement_id INT,
    p_item_ids INT[],
    p_location_id INT,
    p_start_date DATE,
    p_end_date DATE,
    p_payment_method VARCHAR(40)
)
    RETURNS VOID AS
$$
DECLARE
    current_item_id INT;
    available_count INT;
    total_cost      DECIMAL := 0;
    daily_rate      DECIMAL;
    days_rented     INT;
    payment_id      INT;
BEGIN
    -- Рассчитываем количество дней аренды
    days_rented := p_end_date - p_start_date;

    -- Перебираем ID товаров из массива
    FOREACH current_item_id IN ARRAY p_item_ids
        LOOP
            -- Проверяем доступное количество для товара в указанной локации и периоде
            SELECT pa.Number - COALESCE(SUM(CASE
                                                WHEN re.StartDate <= p_end_date AND re.EndDate >= p_start_date THEN 1
                                                ELSE 0
                END), 0)
            INTO available_count
            FROM productAvailability pa
                     LEFT JOIN rentedEquipment re
                               ON pa.InventoryID = re.InventoryID
            WHERE pa.InventoryID = current_item_id
              AND pa.LocationId = p_location_id
            GROUP BY pa.Number;

            -- Если товара недостаточно
            IF available_count <= 0 THEN
                RAISE EXCEPTION 'Not enough items available for InventoryID % at LocationID %', current_item_id, p_location_id;
            END IF;

            -- Добавляем запись об аренде
            INSERT INTO rentedEquipment (RentalAgreementId, InventoryID, StartDate, EndDate)
            VALUES (p_rental_agreement_id, current_item_id, p_start_date, p_end_date);

            -- Увеличиваем общую стоимость аренды
            SELECT i.RentalPricePerDay
            INTO daily_rate
            FROM inventory i
            WHERE i.InventoryID = current_item_id;

            total_cost := total_cost + (daily_rate * days_rented);
        END LOOP;

    -- Добавляем запись о платеже
    INSERT INTO payments (Amount, PaymentDate, PaymentMethod, Status)
    VALUES (total_cost, CURRENT_DATE, p_payment_method, FALSE)
    RETURNING PaymentID INTO payment_id;

    -- Обновляем статус договора и привязываем платеж
    UPDATE rentalAgreements
    SET Status    = TRUE,
        PaymentID = payment_id
    WHERE RentalAgreementID = p_rental_agreement_id;

    RETURN;

EXCEPTION
    WHEN OTHERS THEN
        -- Если возникла ошибка, удаляем договор
        DELETE FROM rentalAgreements WHERE RentalAgreementID = p_rental_agreement_id;
        RAISE; -- Переподнимаем ошибку, чтобы обработать её выше
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_available_periods_by_location_inventory(
    p_location_id INT,
    p_inventory_id INT
)
    RETURNS TABLE
            (
                available_start DATE,
                available_end   DATE
            )
AS
$$
BEGIN
    RETURN QUERY
        WITH booked_periods AS (SELECT re.StartDate, re.EndDate
                                FROM rentedEquipment re
                                         JOIN rentalAgreements ra ON ra.RentalAgreementID = re.RentalAgreementId
                                WHERE re.InventoryID = p_inventory_id
                                  AND EXISTS (SELECT 1
                                              FROM productAvailability pa
                                              WHERE pa.InventoryID = re.InventoryID
                                                AND pa.LocationID = p_location_id)
                                  AND (re.StartDate <= CURRENT_DATE + INTERVAL '1 month'
                                    AND re.EndDate >= CURRENT_DATE)),
             ordered_periods AS (SELECT StartDate,
                                        EndDate,
                                        LAG(EndDate) OVER (ORDER BY StartDate) AS prev_end
                                 FROM booked_periods),
             merged_periods AS (SELECT StartDate,
                                       EndDate,
                                       CASE
                                           WHEN StartDate <= COALESCE(prev_end, CURRENT_DATE - INTERVAL '1 day') THEN 0
                                           ELSE 1
                                           END AS is_new_group
                                FROM ordered_periods),
             grouping_info AS (SELECT StartDate,
                                      EndDate,
                                      SUM(is_new_group) OVER (ORDER BY StartDate) AS group_id
                               FROM merged_periods),
             grouped_periods AS (SELECT MIN(StartDate) AS StartDate, MAX(EndDate) AS EndDate
                                 FROM grouping_info
                                 GROUP BY group_id),
             gaps AS (SELECT COALESCE(
                                             LAG(grouped_periods.EndDate) OVER (ORDER BY grouped_periods.StartDate),
                                             CURRENT_DATE - INTERVAL '1 day'
                             ) + INTERVAL '1 day'                         AS gap_start,
                             grouped_periods.StartDate - INTERVAL '1 day' AS gap_end
                      FROM grouped_periods
                      WHERE grouped_periods.StartDate > CURRENT_DATE

                      UNION ALL

                      -- Добавляем промежуток от конца последнего забронированного периода до конца месяца
                      SELECT MAX(grouped_periods.EndDate) + INTERVAL '1 day' AS gap_start,
                             (CURRENT_DATE + INTERVAL '1 month')::DATE       AS gap_end
                      FROM grouped_periods
                      WHERE NOT EXISTS (SELECT 1
                                        FROM grouped_periods
                                        WHERE grouped_periods.EndDate >= (CURRENT_DATE + INTERVAL '1 month')::DATE)),
             final_gaps AS (SELECT gap_start::DATE AS available_start,
                                   CASE
                                       WHEN gap_end IS NULL THEN (CURRENT_DATE + INTERVAL '1 month')::DATE
                                       ELSE gap_end::DATE
                                       END         AS available_end
                            FROM gaps
                            WHERE gap_start <= COALESCE(gap_end, CURRENT_DATE + INTERVAL '1 month'))
        SELECT final_gaps.available_start,
               final_gaps.available_end
        FROM final_gaps;

    -- Если нет записей в booked_periods, возвращаем весь месяц
    IF NOT FOUND THEN
        RETURN QUERY
            SELECT CURRENT_DATE::DATE, (CURRENT_DATE + INTERVAL '1 month')::DATE;
    END IF;
END;
$$ LANGUAGE plpgsql;



--EQUIPMENT--


--LOCATIONS--

CREATE OR REPLACE FUNCTION add_location(
    p_country VARCHAR(255),
    p_city VARCHAR(255),
    p_address VARCHAR(255),
    p_opening_time TIME,
    p_closing_time TIME,
    p_phone_number VARCHAR(50)
) RETURNS INT AS
$$
DECLARE
    new_location_id INT;
BEGIN
    INSERT INTO locations(country, city, address, openingtime, closingtime, phonenumber)
    VALUES (p_country, p_city, p_address, p_opening_time, p_closing_time, p_phone_number)
    RETURNING LocationID INTO new_location_id;

    RETURN new_location_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_all_locations()
    RETURNS TABLE
            (
                p_location_id  INT,
                p_country      VARCHAR(255),
                p_city         VARCHAR(255),
                p_address      VARCHAR(255),
                p_opening_time TEXT,
                p_closing_time TEXT,
                p_phone_number VARCHAR(50)
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT LocationID                         AS p_location_id,
               Country                            AS p_country,
               City                               AS p_city,
               Address                            AS p_address,
               TO_CHAR(OpeningTime, 'HH24:MI:SS') AS p_opening_time, -- Преобразование времени в строку
               TO_CHAR(ClosingTime, 'HH24:MI:SS') AS p_closing_time, -- Преобразование времени в строку
               PhoneNumber                        AS p_phone_number
        FROM locations
        ORDER BY p_location_id;
END
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION update_location(
    p_location_id INT,
    p_country VARCHAR(255) DEFAULT NULL,
    p_city VARCHAR(255) DEFAULT NULL,
    p_address VARCHAR(255) DEFAULT NULL,
    p_opening_time TIME DEFAULT NULL,
    p_closing_time TIME DEFAULT NULL,
    p_phone_number VARCHAR(50) DEFAULT NULL
) RETURNS BOOLEAN AS
$$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM locations WHERE locationid = p_location_id) THEN
        RAISE EXCEPTION 'Location with ID % does not exist', p_location_id;
    END IF;

    UPDATE locations
    SET country     = COALESCE(p_country, country),
        city        = COALESCE(p_city, city),
        address     = COALESCE(p_address, address),
        openingtime = COALESCE(p_opening_time, openingtime),
        closingtime = COALESCE(p_closing_time, closingtime),
        phonenumber = COALESCE(p_phone_number, phonenumber)
    WHERE locationid = p_location_id;

    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION delete_location(
    p_location_id INT
) RETURNS BOOLEAN AS
$$
BEGIN
    DELETE
    FROM locations
    WHERE LocationID = p_location_id;

    IF FOUND THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$ LANGUAGE plpgsql;
--Триггер на удаление локации?


CREATE OR REPLACE FUNCTION delete_all_locations() RETURNS VOID AS
$$
BEGIN
    DELETE FROM locations;
END;
$$ LANGUAGE plpgsql;

--CASCADE DELETE?


CREATE OR REPLACE FUNCTION create_review(
    p_user_id INT,
    p_item_id INT,
    p_name VARCHAR(255),
    p_rating INT,
    p_comment TEXT,
    p_review_date DATE
)
    RETURNS VOID AS
$$
DECLARE
    rented_id INT;
    review_id INT;
BEGIN
    -- Находим rented_id, который соответствует аренде
    SELECT re.RentedID
    INTO rented_id
    FROM rentedEquipment re
             JOIN rentalAgreements ra ON re.RentalAgreementId = ra.RentalAgreementID
    WHERE ra.CustomerId = p_user_id
      AND re.InventoryID = p_item_id
    LIMIT 1;

    -- Проверяем, что rented_id найден
    IF rented_id IS NULL THEN
        RAISE EXCEPTION 'User % has not rented item %', p_user_id, p_item_id;
    END IF;

    -- Вставляем отзыв и сохраняем его ReviewID в review_id
    INSERT INTO reviews (Name, Rating, Comment, ReviewDate)
    VALUES (p_name, p_rating, p_comment, p_review_date)
    RETURNING ReviewID INTO review_id;

    -- Обновляем rentedEquipment, устанавливая ReviewID для соответствующего rented_id
    UPDATE rentedEquipment
    SET ReviewID = review_id
    WHERE RentedID = rented_id;

    RETURN;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_reviews_by_item_id(p_item_id INT)
    RETURNS TABLE
            (
                review_id   INT,
                name        VARCHAR(255),
                rating      INT,
                comment     TEXT,
                review_date DATE
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT r.ReviewID,
               r.Name,
               r.Rating,
               r.Comment,
               r.ReviewDate
        FROM reviews r
                 JOIN rentedEquipment re ON r.ReviewID = re.ReviewID
        WHERE re.InventoryID = p_item_id;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_reviews_by_user_id(p_user_id INT)
    RETURNS TABLE
            (
                review_id   INT,
                name        VARCHAR(255),
                rating      INT,
                comment     TEXT,
                review_date DATE
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT r.ReviewID,
               r.Name,
               r.Rating,
               r.Comment,
               r.ReviewDate
        FROM reviews r
                 JOIN rentedEquipment re ON r.ReviewID = re.ReviewID
                 JOIN rentalAgreements ra ON re.RentalAgreementId = ra.RentalAgreementID
        WHERE ra.CustomerId = p_user_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_all_reviews()
    RETURNS TABLE
            (
                review_id   INT,
                name        VARCHAR(255),
                rating      INT,
                comment     TEXT,
                review_date DATE,
                item_id     INT
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT r.ReviewID,
               r.Name,
               r.Rating,
               r.Comment,
               r.ReviewDate,
               re.InventoryID AS item_id
        FROM reviews r
                 JOIN rentedEquipment re ON r.ReviewID = re.ReviewID
        ORDER BY rating;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION delete_review_by_id(p_review_id INT)
    RETURNS VOID AS
$$
BEGIN
    -- Удаляем связь с отзывом в rentedEquipment
    UPDATE rentedEquipment
    SET ReviewID = NULL
    WHERE ReviewID = p_review_id;

    -- Удаляем сам отзыв
    DELETE
    FROM reviews
    WHERE ReviewID = p_review_id;
END;
$$ LANGUAGE plpgsql;


--Review--


CREATE OR REPLACE FUNCTION add_product_availability(
    p_inventory_id INT,
    p_location_id INT,
    p_number INT
)
    RETURNS VOID AS
$$
BEGIN
    -- Проверяем, существует ли запись с указанным InventoryID и LocationId
    IF EXISTS (SELECT 1
               FROM productAvailability
               WHERE InventoryID = p_inventory_id
                 AND LocationId = p_location_id) THEN
        -- Обновляем количество для существующей записи
        UPDATE productAvailability
        SET Number = p_number
        WHERE InventoryID = p_inventory_id
          AND LocationId = p_location_id;

        RAISE NOTICE 'Updated existing product availability: InventoryID=%, LocationId=%, Number=%',
            p_inventory_id, p_location_id, p_number;
    ELSE
        -- Добавляем новую запись для нового LocationId
        INSERT INTO productAvailability (InventoryID, LocationId, Number)
        VALUES (p_inventory_id, p_location_id, p_number);

        RAISE NOTICE 'Inserted new product availability: InventoryID=%, LocationId=%, Number=%',
            p_inventory_id, p_location_id, p_number;
    END IF;
END;
$$ LANGUAGE plpgsql;




